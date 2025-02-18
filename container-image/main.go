package main

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/copy"
	_ "github.com/containers/image/v5/docker"
	_ "github.com/containers/image/v5/oci/layout"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports"
	"golang.org/x/sync/errgroup"
)

func getPolicyContext() (*signature.PolicyContext, error) {
	policy := &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}
	return signature.NewPolicyContext(policy)
}

func doImagePull() error {
	ctx := context.Background()
	ociTransport := transports.Get("oci")
	dst, err := ociTransport.ParseReference("my-dir")
	if err != nil {
		return fmt.Errorf("could parse transport reference: %w", err)
	}
	dockerTransport := transports.Get("docker")
	policy, err := getPolicyContext()
	if err != nil {
		return fmt.Errorf("failed to get policy: %w", err)
	}
	images := []string{
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-reflector-controller:v0.33.0",
		"ghcr.io/fluxcd/kustomize-controller:v1.4.0",
		"ghcr.io/fluxcd/notification-controller:v1.4.0",
		"ghcr.io/fluxcd/source-controller:v1.4.1",

		// "ghcr.io/austinabro321/10-layers:v0.0.1", // to test
	}
	for _, image := range images {
		src, err := dockerTransport.ParseReference(fmt.Sprintf("//%s", image))
		if err != nil {
			return fmt.Errorf("couldn't parse: %w", err)
		}
		_, err = copy.Image(ctx, policy, dst, src, nil)
		if err != nil {
			return fmt.Errorf("failed during copy: %w", err)
		}
	}
	return nil
}

func doImagePullConcurrent() error {
	ctx := context.Background()
	ociTransport := transports.Get("oci")
	dst, err := ociTransport.ParseReference("my-dir-concurrent")
	if err != nil {
		return fmt.Errorf("could parse transport reference: %w", err)
	}
	dockerTransport := transports.Get("docker")
	policy, err := getPolicyContext()
	if err != nil {
		return fmt.Errorf("failed to get policy: %w", err)
	}

	images := []string{
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",
		"ghcr.io/fluxcd/image-automation-controller:v0.39.0",		
	}
	eg, ectx := errgroup.WithContext(ctx)
	for _, image := range images {
		image := image
		eg.Go(func() error {
			select {
			case <-ectx.Done():
				return ectx.Err()
			default:
				src, err := dockerTransport.ParseReference(fmt.Sprintf("//%s", image))
				if err != nil {
					return fmt.Errorf("couldn't parse: %w", err)
				}
				_, err = copy.Image(ctx, policy, dst, src, nil)
				if err != nil {
					return fmt.Errorf("failed during copy: %w", err)
				}
			}
			return nil
		})
	}
	return eg.Wait()
}

func main() {
	if err := doImagePull(); err != nil {
		panic(err)
	}
	if err := doImagePullConcurrent(); err != nil {
		panic(err)
	}
}
