rm -rf oras-local
# oras copy --concurrency 10 ghcr.io/austinabro321/one-large-layer:v0.0.1 --to-oci-layout oras-local
oras copy --from-plain-http --concurrency 10 localhost:5000/one-large-layer:v0.0.1 --to-oci-layout oras-local



# oras copy --concurrency 10 nvcr.io/nvidia/k8s/dcgm-exporter:3.3.0-3.2.0-ubuntu22.04 --to-oci-layout oras-local
# oras copy --concurrency 10 nvcr.io/nvidia/k8s/container-toolkit:v1.14.6-ubuntu20.04 --to-oci-layout oras-local
# oras copy --concurrency 10 registry.k8s.io/nfd/node-feature-discovery:v0.14.2 --to-oci-layout oras-local
# oras copy --concurrency 10 nvcr.io/nvidia/gpu-feature-discovery:v0.8.2-ubi8 --to-oci-layout oras-local
# oras copy --concurrency 10 nvcr.io/nvidia/k8s-device-plugin:v0.14.5-ubi8 --to-oci-layout oras-local
# oras copy --concurrency 10 nvcr.io/nvidia/gpu-operator:v23.9.2 --to-oci-layout oras-local
# oras copy --concurrency 10 nvcr.io/nvidia/cloud-native/gpu-operator-validator:v23.9.2 --to-oci-layout oras-local
# oras copy --concurrency 10 nvcr.io/nvidia/cloud-native/k8s-driver-manager:v0.6.5 --to-oci-layout oras-local
# oras copy --concurrency 10 ghcr.io/fluxcd/helm-controller:v1.1.0 --to-oci-layout oras-local
# oras copy --concurrency 10 ghcr.io/fluxcd/image-automation-controller:v0.39.0 --to-oci-layout oras-local
# oras copy --concurrency 10 ghcr.io/fluxcd/image-reflector-controller:v0.33.0 --to-oci-layout oras-local
# oras copy --concurrency 10 ghcr.io/fluxcd/kustomize-controller:v1.4.0 --to-oci-layout oras-local
# oras copy --concurrency 10 ghcr.io/fluxcd/notification-controller:v1.4.0 --to-oci-layout oras-local
# oras copy --concurrency 10 ghcr.io/fluxcd/source-controller:v1.4.1 --to-oci-layout oras-local

# concurrently, I assume is the level of layers pulled at the same time
# real    0m36.959s
# user    0m3.548s
# sys     0m5.779s

# for the zarf.yaml images
# real    2m2.976s
# user    0m9.504s
# sys     0m12.918s