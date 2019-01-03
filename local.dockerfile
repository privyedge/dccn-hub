# Dockerfile.dep is used to build the docker image locally with go dependency,
# in which the RSA key copy part is not needed, the usage description is documented in the README.md

FROM golang:1.11.2-alpine3.8
RUN apk update && \
    apk add --no-cache git && \
    apk add --no-cache --update bash && \
    apk add --no-cache openssh
    
WORKDIR $GOPATH/src/github.com/Ankr-network/dccn-hub/
COPY . $GOPATH/src/github.com/Ankr-network/dccn-hub/

CMD go run cmd/main.go
