FROM alpine:latest



# Create a 4GB random dummy file in a single layer
RUN dd if=/dev/urandom of=/dummyfile bs=1M count=4096