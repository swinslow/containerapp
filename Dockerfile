FROM golang:1.11

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ADD ./webapp /go/src/app

RUN go get -v
