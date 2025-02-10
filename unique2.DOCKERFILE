FROM ghcr.io/fluxcd/image-automation-controller:v0.39.0

COPY dummyfile1 .
COPY dummyfile2 .
COPY dummyfile3 .

CMD [ "echo", "cool2" ]