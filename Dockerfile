FROM golang:alpine3.8
RUN apk update && apk add git && apk add --update bash && apk add openssh

RUN go get github.com/golang/dep/cmd/dep

# Copy the RSA key from the CircleCI build environment into
# this Docker build process in order
# to retrieve the project dependencies from other private repos
# The RSA key must be provisioned into CircleCI ahead of time
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
