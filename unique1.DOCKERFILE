FROM ghcr.io/fluxcd/image-automation-controller:v0.39.0

COPY ./dummy-files/1 .
COPY ./dummy-files/2 .
COPY ./dummy-files/3 .
COPY ./dummy-files/4 .
COPY ./dummy-files/5 .
COPY ./dummy-files/6 .
COPY ./dummy-files/7 .
COPY ./dummy-files/8 .
COPY ./dummy-files/9 .
COPY ./dummy-files/10 .

CMD [ "echo", "cool1" ]