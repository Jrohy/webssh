gox -output="docker/result/webssh_{{.OS}}_{{.Arch}}" -ldflags="-s -w" -os="linux"
cd docker
docker buildx build --platform linux/arm64,linux/amd64,linux/arm,linux/386,linux/ppc64le,linux/s390x -t jrohy/webssh . --push
