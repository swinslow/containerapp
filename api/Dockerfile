FROM golang:1.11

RUN mkdir -p /go/src/github.com/swinslow/containerapp/api
WORKDIR /go/src/github.com/swinslow/containerapp/api

ADD . /go/src/github.com/swinslow/containerapp/api

RUN go get -v ./...
RUN go build
RUN go install github.com/swinslow/containerapp/api