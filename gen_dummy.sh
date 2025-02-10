for i in $(seq 1 10);
do
  dd if=/dev/urandom of="./dummy-files/${i}" bs=1M count=10
done