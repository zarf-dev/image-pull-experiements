# ORAS errors
This error happens concurrently when I was pulling several images at the same time. It was about 50 images, but only five unique images. In fairness we would make images unique in Zarf and I have not been able to repro this error using different images. Also an error like this would only happen with duplicate images at least on the image manifest since it's unique per sha. I'm not sure atm if ORAS has guardrails in place to make sure this doesn't happen with layers


panic: failed to copy: failed to resolve ghcr.io/fluxcd/image-automation-controller:v0.39.0@sha256:48a89734dc82c3a2d4138554b3ad4acf93230f770b3a582f7f48be38436d031c: read failed: sha256:48a89734dc82c3a2d4138554b3ad4acf93230f770b3a582f7f48be38436d031c: application/vnd.oci.image.manifest.v1+json: already exists

# notes
One thing we could theoretically do is to get all the layers we are going to pull and if that layer already exists skip it when pulling the next image. The library would have to support skip functionality, and it might not matter if concurrency always works