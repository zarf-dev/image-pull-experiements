package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
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
	repo := &remote.Repository{}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	dst, err := oci.New(filepath.Join(cwd, "download"))
	if err != nil {
		return err
	}
	// image := "docker.io/library/alpine:latest"
	image := "ghcr.io/fluxcd/image-automation-controller:v0.39.0"
	copyOpts := oras.DefaultCopyOptions
	repo.Reference, err = registry.ParseReference(image)
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
		platformDesc, err := oras.Resolve(ctx, repo, repo.Reference.Reference, resolveOpts)
		if err != nil {
			return err
		}
		image = fmt.Sprintf("%s@%s", image, platformDesc.Digest)
    fmt.Println("new image", image)
	}
	repo.Client = client
	cachePath, err := oci.New(filepath.Join(cwd, "test-cache"))
	if err != nil {
		return err
	}
	cachedDst := cache.New(repo, cachePath)
	desc, err := oras.Copy(ctx, cachedDst, image, dst, "", copyOpts)
	if err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}
	fmt.Println(desc.Digest)
	return nil
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
