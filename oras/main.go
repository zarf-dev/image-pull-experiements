package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/sync/errgroup"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"try-oras.com/cache"
)

// If there is an issue in the cache it will fail, assuming it needs to download the layer that the image is pointed to.
// Would be interesting to potentially make

func doOras() error {
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
	// Can we get this to fail when using unique images
	// images := []string{
	// 	"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
	// 	"ghcr.io/fluxcd/image-reflector-controller:v0.33.0",
	// 	"ghcr.io/fluxcd/kustomize-controller:v1.4.0",
	// 	"ghcr.io/fluxcd/notification-controller:v1.4.0",
	// 	"ghcr.io/fluxcd/source-controller:v1.4.1",

	// 	// "ghcr.io/austinabro321/10-layers:v0.0.1", // to test 
	// }
	images := []string{
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
		"localhost:5000/dummy-unique:v0.0.1",
	}
	copyOpts := oras.DefaultCopyOptions
	eg, ectx := errgroup.WithContext(ctx)
	cachePath, err := oci.New(filepath.Join(cwd, "test-cache"))	
	eg.SetLimit(15)
	for _, image := range images {
		image := image
		eg.Go(func() error {
			select {
			case <-ectx.Done():
				return ectx.Err()
			default:
				localRepo := &remote.Repository{PlainHTTP: true}
				localRepo.Reference, err = registry.ParseReference(image)
				if err != nil {
					return err
				}
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
				fmt.Println("new image", image)
				localRepo.Client = client
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
	return eg.Wait()
}

func main() {
	err := doOras()
	if err != nil {
		panic(err)
	}
}

// cfg, err := config.Load(config.Dir())
// if err != nil {
// 	return err
// }
// configs := []*configfile.ConfigFile{cfg}

// var key = image

// authConf, err := configs[0].GetCredentialsStore(key).Get(key)
// if err != nil {
// 	return fmt.Errorf("unable to get credentials for %s: %w", key, err)
// }

// cred := auth.Credential{
// 	Username:     authConf.Username,
// 	Password:     authConf.Password,
// 	AccessToken:  authConf.RegistryToken,
// 	RefreshToken: authConf.IdentityToken,
// }

// client.Credential = auth.StaticCredential(repo.Reference.Reference, cred)
// repo.Client = client
