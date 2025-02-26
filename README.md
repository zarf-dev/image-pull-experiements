# Image pull experiments

Doing experiments here to see how long it takes to pull images and what the implementations for different libraries look like

docker can create a tarball from an API call directly and it will be in OCI format. However we still need a method to take those files in an OCI directory and move them to the correct OCI directory and update the manifests. There's other aspects to consider such as preserving digests

Next steps:
- See how easy / hard is it to get images from the docker daemon using containers/images
- See how easy / hard it is to get images from the docker daemon natively
- improve how we get auth in oras