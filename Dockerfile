FROM golang:1.16-alpine as builder

RUN echo "Build...."

ENV GO111MODULE=on

WORKDIR /go/src/app/
RUN apk add git
RUN apk add build-base musl-dev

COPY ./ ./
RUN go mod download -x

RUN go get -d -v cmd/api
RUN go build ./cmd/api

ENV gendir=/go/src/app/generated
ENV gendir_test=/go/src/app/test/generated
RUN mkdir -p $gendir
RUN mkdir -p $gendir_test
EXPOSE 8080
ENTRYPOINT ./api