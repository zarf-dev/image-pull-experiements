docker build . -f unique1.DOCKERFILE -t ghcr.io/austinabro321/dummy-unique-1:v0.0.1
docker build . -f unique2.DOCKERFILE -t ghcr.io/austinabro321/dummy-unique-2:v0.0.1
docker build . -f unique3.DOCKERFILE -t ghcr.io/austinabro321/dummy-unique-3:v0.0.1

docker push ghcr.io/austinabro321/dummy-unique-1:v0.0.1
docker push ghcr.io/austinabro321/dummy-unique-2:v0.0.1
docker push ghcr.io/austinabro321/dummy-unique-3:v0.0.1