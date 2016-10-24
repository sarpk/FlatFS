# Pull base image.
FROM ubuntu:latest
RUN apt-get update && \
    apt-get install -y golang && \
    apt-get install -y fuse && \
    apt-get install -y git
ENV GOPATH /usr/gopath/
RUN mkdir /tmp/mountpoint
RUN mkdir /tmp/firstDir
RUN go get github.com/sarpk/FlatFS; exit 0
RUN cd /usr/gopath/src/github.com/sarpk/FlatFS/ && go get .; exit 0
RUN cd /usr/gopath/src/github.com/sarpk/FlatFS/ && go build main.go
