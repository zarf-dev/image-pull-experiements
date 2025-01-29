package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

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
	image := "docker.io/library/alpine:latest"
	copyOpts := oras.DefaultCopyOptions
	repo.Reference, err = registry.ParseReference(image)
	if err != nil {
		return err
	}
	repo.Client = client

	desc, err := oras.Copy(ctx, repo, image, dst, "", copyOpts)
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
