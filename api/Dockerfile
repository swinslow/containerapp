FROM golang:1.11

RUN mkdir -p /go/src/github.com/swinslow/containerapp
WORKDIR /go/src/github.com/swinslow/containerapp

ADD . /go/src/github.com/swinslow/containerapp

RUN go get -v ./...
