FROM --platform=$TARGETPLATFORM alpine

LABEL maintainer="Jrohy <euvkzx@gmail.com>"

ARG TARGETARCH

RUN apk add --no-cache tzdata libc6-compat gcompat

COPY result/webssh_linux_$TARGETARCH /webssh

ENTRYPOINT ["/webssh"]
