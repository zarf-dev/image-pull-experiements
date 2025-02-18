# containers/image

This is the library that skopeo uses under the hood. 

From the testing here it shows that this library does not allow concurrent pulls as the index.json will be overwritten each time. Potentially there are other calls that can be made outside of copy.Image to make this possible, however this will likely be quite complex.