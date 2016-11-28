# Pull base image.
FROM ubuntu:latest
RUN apt-get update && \
    apt-get install -y golang && \
    apt-get install -y fuse && \
    apt-get install -y git
ENV GOPATH /usr/gopath/
RUN mkdir /tmp/mountpoint
RUN mkdir /tmp/flatDir
RUN echo 'Version to be updated fd82660432a2990a9b3af3144cb66a249bae198b'
RUN go get github.com/sarpk/FlatFS
RUN cd /usr/gopath/src/github.com/sarpk/FlatFS/ && go build main.go
