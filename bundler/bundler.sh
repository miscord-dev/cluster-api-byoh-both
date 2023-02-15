#! /bin/bash

set -eux

cd $(dirname $0)

rm -r tmp || true
mkdir -p tmp

ARCH=${ARCH:-amd64}

docker buildx build \
    --platform=linux/$ARCH \
    --build-arg ARCH=$ARCH \
    --build-arg KUBERNETES_VERSION=$KUBERNETES_VERSION \
    --build-arg CONTAINERD_VERSION=$CONTAINERD_VERSION \
    --build-arg KUBERNETES_CNI_VERSION=$KUBERNETES_CNI_VERSION \
    --build-arg KUBERNETES_CRI_VERSION=$KUBERNETES_CRI_VERSION \
    -t intermediate:$$ "$1"
docker run -v $(pwd)/tmp:/host --rm intermediate:$$ bash -c "cp -r /bundler/* /host/"

imgpkg push -i ghcr.io/miscord-dev/byoh-both-bundle:$2 -f ./tmp
