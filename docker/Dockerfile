FROM alpine

LABEL maintainer "Jrohy <euvkzx@gmail.com>"

RUN apk add --no-cache tzdata libc6-compat

COPY webssh /

ENTRYPOINT ["/webssh"]
