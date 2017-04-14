FROM golang:1.8.1
MAINTAINER jennal <jennalcn@gmail.com>
LABEL maintainer "jennalcn@gmail.com"
ADD . /go/src/github.com/jennal/goplay
RUN go install github.com/jennal/goplay