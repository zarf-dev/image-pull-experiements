package main

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/copy"
	_ "github.com/containers/image/v5/docker"
	_ "github.com/containers/image/v5/docker/daemon"
	_ "github.com/containers/image/v5/oci/layout"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports"
	"github.com/containers/image/v5/types"
	"golang.org/x/sync/errgroup"
)

func getPolicyContext() (*signature.PolicyContext, error) {
	policy := &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}
	return signature.NewPolicyContext(policy)
}

func DoImagePull(ociTransport types.ImageTransport) error {
	ctx := context.Background()
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
		"ghcr.io/austinabro321/one-large-layer:v0.0.1",
	}
	for _, image := range images {
		fmt.Println("downloading image", image)
		src, err := dockerTransport.ParseReference(fmt.Sprintf("//%s", image))
		if err != nil {
			return fmt.Errorf("couldn't parse: %w", err)
		}
		_, err = copy.Image(ctx, policy, dst, src, &copy.Options{
			PreserveDigests: true,
			DestinationCtx: &types.SystemContext{
				BlobInfoCacheDir: "local-cache",
			},
		})
		if err != nil {
			return fmt.Errorf("failed during copy: %w", err)
		}
	}
	return nil
}

func DoImagePush(ociTransport types.ImageTransport) error {
	ctx := context.Background()
	srcRef, err := ociTransport.ParseReference("my-dir")
	if err != nil {
		return fmt.Errorf("invalid source name: %v", err)
	}

	dockerTransport := transports.Get("docker")
	destRef, err := dockerTransport.ParseReference("//ghcr.io/austinabro321/dummy-image-1:0.0.1")
	if err != nil {
		return fmt.Errorf("invalid destination name: %v", err)
	}
	policyContext, err := getPolicyContext()
	if err != nil {
		return fmt.Errorf("error loading trust policy: %v", err)
	}
	defer policyContext.Destroy()

	_, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{})
	if err != nil {
		return fmt.Errorf("error copying image: %v", err)
	}

	fmt.Println("Image pushed successfully!")
	return nil
}

func DoImagePullDaemon(ociTransport types.ImageTransport) error {
	ctx := context.Background()
	dst, err := ociTransport.ParseReference("my-dir")
	if err != nil {
		return fmt.Errorf("could parse transport reference: %w", err)
	}
	dockerDaemon := transports.Get("docker-daemon")
	srcRef, err := dockerDaemon.ParseReference("ghcr.io/austinabro321/small-image:1.0.0")
	if err != nil {
		return fmt.Errorf("could parse transport reference: %w", err)
	}
	policy, err := getPolicyContext()
	if err != nil {
		return fmt.Errorf("failed to get policy: %w", err)
	}
	_, err = copy.Image(ctx, policy, dst, srcRef, &copy.Options{})
	if err != nil {
		return fmt.Errorf("failed during copy: %w", err)
	}
	fmt.Println("Image pulled successfully!")
	return nil
}

func main() {
	ociTransport := transports.Get("oci")
	if err := DoImagePull(ociTransport); err != nil {
		panic(err)
	}
	// if err := DoImagePullDaemon(ociTransport); err != nil {
	// 	panic(err)
	// }
	// if err := DoImagePush(ociTransport); err != nil {
	// 	panic(err)
	// }
}

func DoImagePullConcurrent() error {
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
