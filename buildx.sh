#!/bin/bash
LATEST_TAG=`git describe --tags $(git rev-list --tags --max-count=1)`

packr2
gox -output="docker/result/webssh_{{.OS}}_{{.Arch}}" -ldflags="-s -w" -os="linux"
packr2 clean
cd docker
docker buildx build --platform linux/arm64,linux/amd64,linux/arm,linux/386,linux/ppc64le,linux/s390x -t jrohy/webssh:${LATEST_TAG} . --push
docker buildx build --platform linux/arm64,linux/amd64,linux/arm,linux/386,linux/ppc64le,linux/s390x -t jrohy/webssh . --push
docker buildx build --platform linux/arm64 -t jrohy/webssh:arm64 . --push
