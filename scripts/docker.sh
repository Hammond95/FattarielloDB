#!/bin/bash
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
#source "${SCRIPTDIR}/common.sh"
#source "${SCRIPTDIR}/utils.sh"
cd "$SCRIPTDIR"

TAG=$(date +'%Y%m%d%H%M%S')

echo "Login for DockerHub."
read -p "Username: " DOCKER_USERNAME
read -sp "Password: " DOCKER_PASSWORD


echo "Building for Fattariello Server."
docker build \
    -t "fattariellodb-server:$TAG" \
    -f server.Dockerfile \
    --build-arg SERVER_ADDRESS=":8888"\
    ..

docker login \
    -u $DOCKER_USERNAME \
    -p $DOCKER_PASSWORD

docker tag fattariellodb-server:$TAG mdelucadev/fattariellodb-server:$TAG
docker push mdelucadev/fattariellodb-server:$TAG

echo "Building for Fattariello Client."
docker build \
    -t "fattariellodb-client:$TAG" \
    -f client.Dockerfile \
    --build-arg SERVER_ADDRESS=":8888"\
    ..

docker login \
    -u $DOCKER_USERNAME \
    -p $DOCKER_PASSWORD

docker tag fattariellodb-client:$TAG mdelucadev/fattariellodb-client:$TAG
docker push mdelucadev/fattariellodb-client:$TAG
