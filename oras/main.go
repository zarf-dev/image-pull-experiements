package main

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/docker/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/sync/errgroup"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"try-oras.com/cache"
)

// If there is an issue in the cache it will fail, assuming it needs to download the layer that the image is pointed to.
// Would be interesting to potentially make

// images := []string{
// "ghcr.io/austinabro321/10-layers:v0.0.1",
// 	"ghcr.io/austinabro321/dummy-image-1:0.0.1",
// 	"ghcr.io/austinabro321/dummy-image-2:0.0.1",
// 	"ghcr.io/austinabro321/dummy-image-3:0.0.1",
// 	"ghcr.io/austinabro321/dummy-image-4:0.0.1",
// 	"ghcr.io/austinabro321/dummy-image-5:0.0.1",
// 	"ghcr.io/austinabro321/dummy-image-6:0.0.1",
// 	"ghcr.io/austinabro321/dummy-image-7:0.0.1",
// 	"ghcr.io/austinabro321/dummy-image-8:0.0.1",
// 	"ghcr.io/austinabro321/dummy-image-9:0.0.1",
// 	"ghcr.io/austinabro321/dummy-image-10:0.0.1",
// }

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func doOrasPullConcurrent() error {
	start := time.Now()
	ctx := context.Background()
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	dst, err := oci.New(filepath.Join(cwd, "download"))
	if err != nil {
		return err
	}
	// Can we get this to fail when using unique images
	// images := []string{
	// 	"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
	// 	"ghcr.io/fluxcd/image-reflector-controller:v0.33.0",
	// 	"ghcr.io/fluxcd/kustomize-controller:v1.4.0",
	// 	"ghcr.io/fluxcd/notification-controller:v1.4.0",
	// 	"ghcr.io/fluxcd/source-controller:v1.4.1",
	// }
	images := []string{
		// "ghcr.io/austinabro321/10-layers:v0.0.1",
		"ghcr.io/austinabro321/dummy-image-1:0.0.1",
		"ghcr.io/austinabro321/dummy-image-2:0.0.1",
		"ghcr.io/austinabro321/dummy-image-3:0.0.1",
		"ghcr.io/austinabro321/dummy-image-4:0.0.1",
		"ghcr.io/austinabro321/dummy-image-5:0.0.1",
		"ghcr.io/austinabro321/dummy-image-6:0.0.1",
		"ghcr.io/austinabro321/dummy-image-7:0.0.1",
		"ghcr.io/austinabro321/dummy-image-8:0.0.1",
		"ghcr.io/austinabro321/dummy-image-9:0.0.1",
		"ghcr.io/austinabro321/dummy-image-10:0.0.1",
	}
	platform := ocispec.Platform{
		Architecture: "amd64",
		OS:           "linux",
	}

	// egCtx :=
	g, _ := errgroup.WithContext(ctx)
	var (
		mu            sync.Mutex
		layersToPull  []string
		imagesWLayers = make(map[string][]string)
	)

	for _, image := range images {
		// capture 'image' in local var for the goroutine
		img := image

		g.Go(func() error {
			localRepo := &remote.Repository{PlainHTTP: true}
			var err error

			// Parse reference
			localRepo.Reference, err = registry.ParseReference(img)
			if err != nil {
				return err
			}

			// Grab credentials
			creds, err := getCreds(localRepo)
			if err != nil {
				return err
			}

			// Set up client with credentials
			client := auth.DefaultClient
			client.Credential = creds
			localRepo.Client = client

			// Resolve the reference
			resolveOpts := oras.ResolveOptions{
				TargetPlatform: &platform,
			}
			platformDesc, err := oras.Resolve(ctx, localRepo, localRepo.Reference.Reference, resolveOpts)
			if err != nil {
				return err
			}

			// Fetch the manifest contents
			b, err := content.FetchAll(ctx, localRepo, platformDesc)
			if err != nil {
				return err
			}

			var manifest ocispec.Manifest
			if err := json.Unmarshal(b, &manifest); err != nil {
				return err
			}

			// Figure out which layers are new
			necessaryLayers := []string{}
			for _, layer := range manifest.Layers {
				layerDigest := layer.Digest.String()

				// Lock before checking and modifying shared slice
				mu.Lock()
				if !contains(layersToPull, layerDigest) {
					layersToPull = append(layersToPull, layerDigest)
					necessaryLayers = append(necessaryLayers, layerDigest)
				}
				mu.Unlock()
			}

			// Store the new layers for this image
			mu.Lock()
			imagesWLayers[img] = necessaryLayers
			mu.Unlock()

			return nil
		})
	}
	g.Wait()
	eg, ectx := errgroup.WithContext(ctx)
	cachePath, err := oci.NewWithContext(ctx, filepath.Join(cwd, "test-cache"))
	eg.SetLimit(10)
	for image, neededLayers := range imagesWLayers {
		image := image
		neededLayers := neededLayers
		eg.Go(func() error {
			select {
			case <-ectx.Done():
				return ectx.Err()
			default:
				copyOpts := oras.DefaultCopyOptions
				copyOpts.Concurrency = 10
				copyOpts.PreCopy = func(ctx context.Context, src ocispec.Descriptor) error {
					if src.MediaType == ocispec.MediaTypeImageLayer || src.MediaType == ocispec.MediaTypeImageLayerGzip || src.MediaType == ocispec.MediaTypeImageLayerZstd {
						if !contains(neededLayers, src.Digest.String()) {
							// fmt.Println("skipping layer", src.Digest.String())
							return oras.SkipNode
						}
					}
					return nil
				}
				copyOpts.WithTargetPlatform(&platform)
				localRepo := &remote.Repository{PlainHTTP: true}
				localRepo.Reference, err = registry.ParseReference(image)
				if err != nil {
					return err
				}
				client := auth.DefaultClient
				creds, err := getCreds(localRepo)
				if err != nil {
					return err
				}
				client.Credential = creds
				localRepo.Client = client
				fmt.Println("downloading image", image)
				cachedDst := cache.New(localRepo, cachePath)
				desc, err := oras.Copy(ctx, cachedDst, image, dst, "", copyOpts)
				if err != nil {
					return fmt.Errorf("failed to copy: %w", err)
				}
				fmt.Println("finished copying image", desc.Digest)
				return nil
			}
		})
	}
	err = eg.Wait()
	if err != nil {
		return err
	}
	fmt.Println("finished concurrent pulling in", time.Since(start))
	return nil
}

func extractTar(tarPath, destDir string) error {
	f, err := os.Open(tarPath)
	if err != nil {
		return fmt.Errorf("failed to open tar file: %w", err)
	}
	defer f.Close()

	tr := tar.NewReader(f)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			// End of archive
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		targetPath := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory if it doesn't exist
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}

		case tar.TypeReg:
			// Ensure parent directory is created
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", targetPath, err)
			}

			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file data for %s: %w", targetPath, err)
			}
			outFile.Close()
		}
	}

	return nil
}

func main() {
	var err error
	err = doOrasPullConcurrent()
	if err != nil {
		panic(err)
	}
	// err = DoOrasPush()
	// if err != nil {
	// 	panic(err)
	// }
	err = PullFromDocker()
	if err != nil {
		panic(err)
	}
}

// PullFromDocker pulls a container image from the Docker daemon and adds it to an OCI-format directory.
func PullFromDocker() error {
	ctx := context.Background()
	imageName := "ghcr.io/local/small-image:0.0.1"

	// Initialize Docker client
	// TODO add things like Host here
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()	
	// Save the image to a tar stream
	// TODO set platform as option during save
	imageReader, err := cli.ImageSave(ctx, []string{imageName})
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}
	defer imageReader.Close()

	tarFile, err := os.Create("image.tar")
	if err != nil {
		return fmt.Errorf("failed to create tar file: %w", err)
	}
	defer tarFile.Close()

	// Read bytes from imageReader and write them to tarFile
	if _, err := io.Copy(tarFile, imageReader); err != nil {
		return fmt.Errorf("error writing image to tar file: %w", err)
	}

	err = extractTar("image.tar", "docker-image")
	if err != nil {
		return err
	}

	b, err := os.ReadFile(filepath.Join("docker-image", "index.json"))
	if err != nil {
		return fmt.Errorf("failed to read index.json: %w", err)
	}
	var index ocispec.Index
	if err := json.Unmarshal(b, &index); err != nil {
		return fmt.Errorf("failed to unmarshal index.json: %w", err)
	}
	// Indexes should always contain exactly one manifests for the single image we are pulling
	if len(index.Manifests) != 1 {
		return fmt.Errorf("index.json does not contain one manifest")
	}
	// Docker does not properly set the image name annotation, we set it here so that ORAS can pick up the image 
	index.Manifests[0].Annotations[ocispec.AnnotationRefName] = imageName
	b, err = json.Marshal(index)
	if err != nil {
		return fmt.Errorf("failed to marshal index.json: %w", err)
	}
	err = os.WriteFile(filepath.Join("docker-image", "index.json"), b, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write index.json: %w", err)
	}

	dockerImageSrc, err := oci.New("docker-image")
	if err != nil {
		return fmt.Errorf("failed to create OCI store: %w", err)
	}

	ociDst, err := oci.New("download")
	if err != nil {
		return err
	}

	// Import the image into OCI store
	fmt.Printf("Importing image %s into OCI directory %s...\n", imageName, "download")
	desc, err := oras.Copy(ctx, dockerImageSrc, imageName, ociDst, "", oras.DefaultCopyOptions)
	if err != nil {
		return fmt.Errorf("failed to import image into OCI store: %w", err)
	}

	fmt.Printf("Successfully imported image %s to OCI directory\n", desc.Digest)
	return nil
}

func DoOrasPush() error {
	ctx := context.Background()

	localStore, err := oci.New("download")
	if err != nil {
		return err
	}

	reference := "ghcr.io/austinabro321/dummy-image-1:0.0.1"
	repo, err := remote.NewRepository(reference)
	if err != nil {
		return fmt.Errorf("failed to create repo: %w", err)
	}

	client := auth.DefaultClient
	credFunc, err := getCreds(repo)
	if err != nil {
		return err
	}
	client.Credential = credFunc
	repo.Client = client
	oras.DefaultCopyOptions.Concurrency = 10

	desc, err := oras.Copy(ctx, localStore, reference, repo, "", oras.DefaultCopyOptions)
	if err != nil {
		return err
	}

	fmt.Printf("Pushed %s with digest %s\n", reference, desc.Digest)

	return nil
}

func DoOrasPull() error {
	// Potentially use this as a way to set creds
	// credentials.NewStoreFromDocker()
	ctx := context.Background()
	client := auth.DefaultClient
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	dst, err := oci.New(filepath.Join(cwd, "download"))
	if err != nil {
		return err
	}
	images := []string{
		"ghcr.io/austinabro321/dummy-image-1:0.0.1",
		"ghcr.io/fluxcd/image-reflector-controller:v0.33.0",
		"ghcr.io/fluxcd/kustomize-controller:v1.4.0",
		"ghcr.io/fluxcd/notification-controller:v1.4.0",
		"ghcr.io/fluxcd/source-controller:v1.4.1",
	}
	copyOpts := oras.DefaultCopyOptions
	cachePath, err := oci.New(filepath.Join(cwd, "test-cache"))
	if err != nil {
		return err
	}
	for _, image := range images {
		localRepo := &remote.Repository{PlainHTTP: true}
		localRepo.Reference, err = registry.ParseReference(image)
		if err != nil {
			return err
		}
		creds, err := getCreds(localRepo)
		if err != nil {
			return fmt.Errorf("failed to get credentials: %w", err)
		}
		client.Credential = creds
		localRepo.Client = client
		if !strings.Contains(image, "@") {
			platform := ocispec.Platform{
				Architecture: "amd64",
				OS:           "linux",
			}
			resolveOpts := oras.ResolveOptions{
				TargetPlatform: &platform,
			}
			platformDesc, err := oras.Resolve(ctx, localRepo, localRepo.Reference.Reference, resolveOpts)
			if err != nil {
				return err
			}
			image = fmt.Sprintf("%s@%s", image, platformDesc.Digest)
		}
		cachedDst := cache.New(localRepo, cachePath)
		desc, err := oras.Copy(ctx, cachedDst, image, dst, "", copyOpts)
		if err != nil {
			return fmt.Errorf("failed to copy: %w", err)
		}
		fmt.Println("finished copying image", desc.Digest)
	}
	return nil
}

func getCreds(localRepo *remote.Repository) (auth.CredentialFunc, error) {
	cfg, err := config.Load(config.Dir())
	if err != nil {
		return nil, err
	}
	configs := []*configfile.ConfigFile{cfg}
	key := localRepo.Reference.Registry
	if key == "registry-1.docker.io" {
		// Docker stores its credentials under the following key, otherwise credentials use the registry URL
		key = "https://index.docker.io/v1/"
	}

	authConf, err := configs[0].GetCredentialsStore(key).Get(key)
	if err != nil {
		return nil, fmt.Errorf("unable to get credentials for %s: %w", key, err)
	}

	cred := auth.Credential{
		Username:     authConf.Username,
		Password:     authConf.Password,
		AccessToken:  authConf.RegistryToken,
		RefreshToken: authConf.IdentityToken,
	}
	return auth.StaticCredential(localRepo.Reference.Registry, cred), nil
}
