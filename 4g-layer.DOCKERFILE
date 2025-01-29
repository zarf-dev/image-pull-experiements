FROM alpine:latest

# Create a 4GB dummy file
RUN dd if=/dev/zero of=/dummyfile bs=1M count=4096
