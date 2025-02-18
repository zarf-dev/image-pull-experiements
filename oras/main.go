package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
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
	var layersToPull []string
	imagesWLayers := map[string][]string{}
	imageShas := []string{}
	for _, image := range images {
		if strings.Contains(image, "@") {
			imageShas = append(imageShas, image)
			continue
		}
		localRepo := &remote.Repository{PlainHTTP: true}
		localRepo.Reference, err = registry.ParseReference(image)
		if err != nil {
			return err
		}
		creds, err := getCreds(localRepo)
		if err != nil {
			return err
		}
		client := auth.DefaultClient
		client.Credential = creds
		localRepo.Client = client
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
		b, err := content.FetchAll(ctx, localRepo, platformDesc)
		if err != nil {
			return err
		}
		var manifest ocispec.Manifest
		err = json.Unmarshal(b, &manifest)
		if err != nil {
			return err
		}
		necessaryLayers := []string{}
		for _, layer := range manifest.Layers {
			if !contains(layersToPull, layer.Digest.String()) {
				layersToPull = append(layersToPull, layer.Digest.String())
				necessaryLayers = append(necessaryLayers, layer.Digest.String())
			}
		}
		imagesWLayers[image] = necessaryLayers
		imageShas = append(imageShas, fmt.Sprintf("%s@%s", image, platformDesc.Digest))
	}
	fmt.Println(imageShas)
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
				copyOpts.PreCopy = func(ctx context.Context, src ocispec.Descriptor) error {
					if src.MediaType == ocispec.MediaTypeImageLayer || src.MediaType == ocispec.MediaTypeImageLayerGzip || src.MediaType == ocispec.MediaTypeImageLayerZstd {
						if !contains(neededLayers, src.Digest.String()) {
							fmt.Println("skipping layer", src.Digest.String())
							return oras.SkipNode
						}
					}
					return nil
				}
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

func main() {
	var err error
	err = doOrasPullConcurrent()
	if err != nil {
		panic(err)
	}
	// err = DoOrasPull()
	// if err != nil {
	// 	panic(err)
	// }
}

func DoOrasPull() error {
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
