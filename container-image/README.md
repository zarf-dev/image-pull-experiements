# containers/image

This is the library that skopeo uses under the hood. 

The dealbreaker in this library is that it does not have a blob cache. There is a blobinfo cache with metadata in it, but it does not allow caching the actual blobs. This is too core of a functionality to Zarf that causes us to exclude it. 

From the testing here it shows that this library does not allow concurrent pulls as the index.json will be overwritten each time. Potentially there are other calls that can be made outside of copy.Image to make this possible, however this will likely be quite complex.

When grabbing an image from the docker daemon the sha of the layers is different than from crane. For the manifest this makes sense because the media type is now using the oci type. I'm not sure why it happens with other types. Looks like there is a preserve digest options, we may want to do that. We also need to check if this is something that can happen with non docker images / images pull from a manifest. I should also check if this is something that currently happens happen with crane / docker by default