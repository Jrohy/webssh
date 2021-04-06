#!/bin/bash
VERSION=`git describe --tags $(git rev-list --tags --max-count=1)`
NOW=`TZ=Asia/Shanghai date "+%Y%m%d-%H%M"`
GO_VERSION=`go version|awk '{print $3,$4}'`
GIT_VERSION=`git rev-parse HEAD`
LDFLAGS="-w -s -X 'main.version=$VERSION' -X 'main.buildDate=$NOW' -X 'main.goVersion=$GO_VERSION' -X 'main.gitVersion=$GIT_VERSION'"

gox -output="docker/result/webssh_{{.OS}}_{{.Arch}}" -ldflags="$LDFLAGS" -os="linux"

cd docker
docker buildx build --platform linux/arm64,linux/amd64,linux/arm,linux/386,linux/ppc64le,linux/s390x -t jrohy/webssh:${VERSION} . --push
docker buildx build --platform linux/arm64,linux/amd64,linux/arm,linux/386,linux/ppc64le,linux/s390x -t jrohy/webssh . --push
docker buildx build --platform linux/arm64 -t jrohy/webssh:arm64 . --push
