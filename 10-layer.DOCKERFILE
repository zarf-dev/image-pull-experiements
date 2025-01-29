FROM alpine:latest

RUN dd if=/dev/urandom of=/dummyfile1 bs=1M count=400

RUN dd if=/dev/urandom of=/dummyfile2 bs=1M count=400

RUN dd if=/dev/urandom of=/dummyfile3 bs=1M count=400

RUN dd if=/dev/urandom of=/dummyfile4 bs=1M count=400

RUN dd if=/dev/urandom of=/dummyfile5 bs=1M count=400

RUN dd if=/dev/urandom of=/dummyfile6 bs=1M count=400

RUN dd if=/dev/urandom of=/dummyfile7 bs=1M count=400

RUN dd if=/dev/urandom of=/dummyfile8 bs=1M count=400

RUN dd if=/dev/urandom of=/dummyfile9 bs=1M count=400

RUN dd if=/dev/urandom of=/dummyfile10 bs=1M count=400