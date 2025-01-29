rm -rf skopeo-local

# skopeo sync --src docker --dest dir nvcr.io/nvidia/k8s/container-toolkit:v1.14.6-ubuntu20.04 skopeo-local
# skopeo sync --src docker --dest dir nvcr.io/nvidia/k8s/dcgm-exporter:3.3.0-3.2.0-ubuntu22.04 skopeo-local
skopeo sync --src-tls-verify=false --src yaml --dest dir skopeo.yaml skopeo-local
# skopeo copy --src-tls-verify=false docker://localhost:5000/one-large-layer:v0.0.1 dir:skopeo-local

# for ghcr.io/austinabro321/one-large-layer:v0.0.1
# real    2m0.695s
# user    0m11.054s
# sys     0m14.416s