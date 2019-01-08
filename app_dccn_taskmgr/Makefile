.PHONY: proto build test run deploy

GOPATH:=$(shell go env GOPATH)

all: run test deploy devtools clean

proto:
	@protoc -I. \
	    -I${GOPATH}/src:. \
	    --micro_out=. \
	    --go_out=. \
	    proto/taskmgr/taskmgr.proto

build: proto
	docker build . -t app_dccn_taskmgr:latest

run: build
	docker run -d \
		-p 50051 \
		-e DB_HOST="127.0.0.1" \
		-e DB=dccn \
		app_dccn_taskmgr

clean:
	rm app_dccn_taskmgr

test:
	go test -v ./... -cover -race

deploy:
	@echo "docker push"

devtools:
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go get github.com/micro/protoc-gen-micro
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
	@type "protoc-gen-micro" 2> /dev/null || echo 'Please install protoc-gen-micro'

