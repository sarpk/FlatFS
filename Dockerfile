# Pull base image.
FROM ubuntu:latest
RUN apt-get update && \
    apt-get install -y golang && \
    apt-get install -y fuse && \
    apt-get install -y git
ENV GOPATH /usr/gopath/
RUN mkdir /tmp/mountpoint
RUN go get github.com/hanwen/go-fuse; exit 0