
GOPATH:=$(shell go env GOPATH)

all: run test deploy devtools local

SUBDIRS := $(wildcard proto/**/*.proto)
.PHONY: proto build test run deploy $(SUBDIRS) clean

proto: $(SUBDIRS)
$(SUBDIRS):
	@protoc -I. \
	    -I${GOPATH}/src:. \
	    --micro_out=. \
	    --go_out=. \
		$@

build: proto
	docker build . -t gateway:latest

run: build
	docker run -d \
		-p 8080:8080\
		microhq/micro api --handler=rpc \
		-e MICRO_API_ADDRESS=:8080
		gateway

local:
	echo "ensure run micro api --handler=rpc"
	go run main.go

clean:
	rm gateway

test:
	go test -v ./... -cover -race

deploy:
	@echo "docker push"

devtools:
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go get github.com/micro/protoc-gen-micro
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
	@type "protoc-gen-micro" 2> /dev/null || echo 'Please install protoc-gen-micro'

