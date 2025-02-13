package main

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/copy"
	_ "github.com/containers/image/v5/docker"
	_ "github.com/containers/image/v5/oci/layout"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports"
)

func getPolicyContext() (*signature.PolicyContext, error) {
	var policy *signature.Policy // This could be cached across calls in opts.
	policy = &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}
	// if err != nil {
	// 	return nil, err
	// }
	// policy, err = signature.NewPolicyFromFile(opts.policyPath)
	// if err != nil {
	// 	return nil, err
	// }
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
	src, err := dockerTransport.ParseReference("//ghcr.io/fluxcd/image-automation-controller:v0.39.0")
	if err != nil {
		return fmt.Errorf("couldn't parse: %w", err)
	}
	_, err = copy.Image(ctx, policy, dst, src, nil)
	if err != nil {
		return fmt.Errorf("failed during copy: %w", err)
	}
	return nil
}

func main() {
	fmt.Println("hello world")
	if err := doImagePull(); err != nil {
		panic(err)
	}
}
