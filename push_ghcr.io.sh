BUILD_DIR="./docker_builds"
for i in $(seq 1 25); do
    DOCKERFILE="$BUILD_DIR/Dockerfile_$i"
    IMAGE_TAG="ghcr.io/austinabro321/dummy-image-$i:0.0.1"
    docker build -f "$DOCKERFILE" -t "$IMAGE_TAG" .
    docker push $IMAGE_TAG
done