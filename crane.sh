rm -rf crane-local
crane pull nvcr.io/nvidia/k8s/dcgm-exporter:3.3.0-3.2.0-ubuntu22.04 crane-local --format=oci
crane pull nvcr.io/nvidia/k8s/container-toolkit:v1.14.6-ubuntu20.04 crane-local --format=oci
crane pull registry.k8s.io/nfd/node-feature-discovery:v0.14.2 crane-local --format=oci
crane pull nvcr.io/nvidia/gpu-feature-discovery:v0.8.2-ubi8 crane-local --format=oci
crane pull nvcr.io/nvidia/k8s-device-plugin:v0.14.5-ubi8 crane-local --format=oci
crane pull nvcr.io/nvidia/gpu-operator:v23.9.2 crane-local --format=oci
crane pull nvcr.io/nvidia/cloud-native/gpu-operator-validator:v23.9.2 crane-local --format=oci
crane pull nvcr.io/nvidia/cloud-native/k8s-driver-manager:v0.6.5 crane-local --format=oci
crane pull ghcr.io/fluxcd/helm-controller:v1.1.0 crane-local --format=oci
crane pull ghcr.io/fluxcd/image-automation-controller:v0.39.0 crane-local --format=oci
crane pull ghcr.io/fluxcd/image-reflector-controller:v0.33.0 crane-local --format=oci
crane pull ghcr.io/fluxcd/kustomize-controller:v1.4.0 crane-local --format=oci
crane pull ghcr.io/fluxcd/notification-controller:v1.4.0 crane-local --format=oci
crane pull ghcr.io/fluxcd/source-controller:v1.4.1 crane-local --format=oci 

# for nvcr.io/nvidia/k8s/container-toolkit:v1.14.6-ubuntu20.04
# real    0m33.845s
# user    0m3.334s
# sys     0m4.945s

# for the zarf.yaml images
# real    5m47.000s
# user    0m9.282s
# sys     0m13.437s