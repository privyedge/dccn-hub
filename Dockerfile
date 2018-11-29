FROM golang:alpine as builder
RUN apk update && apk add git && apk add --update bash && apk add openssh


COPY test.go test.go


CMD go run test.go
