#!/bin/bash
version=`git describe --tags $(git rev-list --tags --max-count=1)`
now=`TZ=Asia/Shanghai date "+%Y%m%d-%H%M"`
go_version=`go version|awk '{print $3,$4}'`
git_version=`git rev-parse HEAD`
ldflags="-w -s -X 'main.version=$version' -X 'main.buildDate=$now' -X 'main.goVersion=$go_version' -X 'main.gitVersion=$git_version'"

gox -output="docker/result/webssh_{{.OS}}_{{.Arch}}" -ldflags="$ldflags" -os="linux"

cd docker
docker buildx build --platform linux/arm64,linux/amd64,linux/arm,linux/386,linux/ppc64le,linux/s390x -t jrohy/webssh:${version} . --push
docker buildx build --platform linux/arm64,linux/amd64,linux/arm,linux/386,linux/ppc64le,linux/s390x -t jrohy/webssh . --push
docker buildx build --platform linux/arm64 -t jrohy/webssh:arm64 . --push
