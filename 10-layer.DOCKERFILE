FROM alpine:latest

RUN dd if=/dev/zero of=/dummyfile1 bs=1M count=400

RUN dd if=/dev/zero of=/dummyfile2 bs=1M count=400

RUN dd if=/dev/zero of=/dummyfile3 bs=1M count=400

RUN dd if=/dev/zero of=/dummyfile4 bs=1M count=400

RUN dd if=/dev/zero of=/dummyfile5 bs=1M count=400

RUN dd if=/dev/zero of=/dummyfile6 bs=1M count=400

RUN dd if=/dev/zero of=/dummyfile7 bs=1M count=400

RUN dd if=/dev/zero of=/dummyfile8 bs=1M count=400

RUN dd if=/dev/zero of=/dummyfile9 bs=1M count=400

RUN dd if=/dev/zero of=/dummyfile10 bs=1M count=400