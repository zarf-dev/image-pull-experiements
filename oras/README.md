# Concurrency support
ORAS seems to support concurrency for the oci store, which holds images in the oci-format. https://github.com/oras-project/oras-go/blob/8d44f2342f185e0195acbdccc30eb9ac2d741d20/content/oci/oci.go#L78

## ORAS errors
This error happens concurrently when I was pulling several images at the same time. It was about 50 images, but only five unique images. In fairness we would make images unique in Zarf and I have not been able to repro this error using different images. Also an error like this would only happen with duplicate images at least on the image manifest since it's unique per sha.

panic: failed to copy: failed to resolve ghcr.io/fluxcd/image-automation-controller:v0.39.0@sha256:48a89734dc82c3a2d4138554b3ad4acf93230f770b3a582f7f48be38436d031c: read failed: sha256:48a89734dc82c3a2d4138554b3ad4acf93230f770b3a582f7f48be38436d031c: application/vnd.oci.image.manifest.v1+json: already exists

# notes
Another speed up unique to ORAS is the ability to skip a blob that we don't need. To do this we find the layers we are going to pull across all images. Any duplicates are pulled only once. This won't matter for layers already in the cache, but it ensures two images don't try to pull at the same time. A new image could come into the queue with the same layers as an image that is about to finish. No reason in this case for the image to re-pull that layer. This also removes the opportunity for any file system issues, oras probably solves this already with syncs, but this could lower the amount of waiting depending on how the sync is working. 

Docker seems to work, you just need to grab it, unarchive it as an OCI directory, and add the right annotation