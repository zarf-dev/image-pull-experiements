rm -rf skopeo-local

# skopeo sync --src docker --dest dir nvcr.io/nvidia/k8s/container-toolkit:v1.14.6-ubuntu20.04 skopeo-local
# skopeo sync --src docker --dest dir nvcr.io/nvidia/k8s/dcgm-exporter:3.3.0-3.2.0-ubuntu22.04 skopeo-local
skopeo sync --src yaml --dest dir skopeo.yaml skopeo-local

# for nvcr.io/nvidia/k8s/container-toolkit:v1.14.6-ubuntu20.04
# real    0m20.678s
# user    0m2.347s
# sys     0m3.508s

# for the entire yaml
# real    1m13.286s
# user    0m6.504s
# sys     0m10.143s