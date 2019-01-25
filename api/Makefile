.PHONY: build test run deploy local clean

GOPATH:=$(shell go env GOPATH)

all: run test deploy devtools

build:
	docker build . -t dccn-api:latest

run: build
	docker run -d \
		-p 50051 \
		$(MICRO_ENV) \
		app_dccn_dcmgr

local:
	$(MICRO_ENV) \
	go run main.go

clean:
	rm api

test:
	go test -v ./... -cover -race

deploy:
	@echo "docker push"

devtools:
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go get github.com/micro/protoc-gen-micro
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
	@type "protoc-gen-micro" 2> /dev/null || echo 'Please install protoc-gen-micro'

define MICRO_ENV
	MICRO_REGISTRY="mdns" \
	MICRO_SERVER_ADDRESS="127.0.0.1:50051" \
	MICRO_SERVER_VERSION="v1.0" \
	MICRO_REGISTER_TTL=15 \
	MICRO_REGISTER_INTERVAL=15
endef

