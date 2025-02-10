# Image pull experiments

Doing experiments here to see how long it takes to pull images and what the implementations for different libraries look like

One thing we could theoretically do is to get all the layers we are going to pull and if that layer already exists skip it when pulling the next image. The library would have to support skip functionality, and it might not matter if concurrency always works

Next steps:
- determine cost of keeping both crane and oras in the repo. Crane would be only used to pull images from the daemon
- Determine 100% that ORAS cannot pull images from the docker client
- Determine 100% that containers/image is not able to export to oci format