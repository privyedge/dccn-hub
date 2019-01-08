.PHONY: proto build test run deploy local clean

GOPATH:=$(shell go env GOPATH)

all: run test deploy devtools

proto:
	@protoc -I. \
	    -I${GOPATH}/src:. \
	    --micro_out=. \
	    --go_out=. \
	    proto/dcmgr/dcmgr.proto

build: proto
	docker build . -t app_dccn_dcmgr:latest

run: build
	docker run -d \
		-p 50051 \
		app_dccn_dcmgr

local:
		DB_HOST="127.0.0.1" \
		DB=dccn \
		COLLECTION=datacenter \
		POOL_LIMIT=4096 \
		TIMEOUT=15 \
		MICRO_REGISTER_TTL=30 \
		MICRO_REGISTER_INTERVAL=30 \
		MICRO_SERVER_NAME=go.micro.srv.dcmgr \
		MICRO_SERVER_VERSION=v0.1 \
		MICRO_SERVER_ADDRESS=:50051 \
		go run main.go

clean:
	rm app_dccn_dcmgr

test:
	go test -v ./... -cover -race

deploy:
	@echo "docker push"

devtools:
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go get github.com/micro/protoc-gen-micro
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
	@type "protoc-gen-micro" 2> /dev/null || echo 'Please install protoc-gen-micro'

