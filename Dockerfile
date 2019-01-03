# To do: Use multi-stage docker build process
FROM golang:1.11.2-alpine3.8

RUN apk update && \
    apk add --no-cache git && \
    apk add --no-cache --update bash && \
    apk add --no-cache openssh

RUN go get github.com/golang/dep/cmd/dep

COPY id_rsa /root/.ssh/
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN chmod go-w /root
RUN chmod 700 /root/.ssh
RUN chmod 600 /root/.ssh/id_rsa

WORKDIR $GOPATH/src/github.com/Ankr-network/dccn-hub/
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only
COPY . $GOPATH/src/github.com/Ankr-network/dccn-hub/

CMD go run cmd/main.go mongo
