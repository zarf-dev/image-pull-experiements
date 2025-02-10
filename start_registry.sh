docker run -p 5000:5000 -d distribution/distribution:2.8.3
docker build . -f unique1.DOCKERFILE -t localhost:5000/dummy-unique:v0.0.1
docker push localhost:5000/dummy-unique:v0.0.1