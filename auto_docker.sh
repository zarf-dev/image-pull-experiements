# Base image and common layers
BASE_IMAGE="ghcr.io/fluxcd/image-automation-controller:v0.39.0"
DUMMY_FILES=(1 2 3 4 5 6 7 8 9 10)

# Directory to store temporary Dockerfiles
BUILD_DIR="./docker_builds"
mkdir -p "$BUILD_DIR"

for i in $(seq 1 25); do
    DOCKERFILE="$BUILD_DIR/Dockerfile_$i"
    IMAGE_TAG="myimage:tag_$i"

    # Generate the Dockerfile
    {
        echo "FROM $BASE_IMAGE"
        for file in "${DUMMY_FILES[@]}"; do
            echo "COPY ./dummy-files/$file ."
        done
        echo "CMD [\"echo\", \"cool$i\"]"
    } > "$DOCKERFILE"

    echo "created file $IMAGE_TAG"
done