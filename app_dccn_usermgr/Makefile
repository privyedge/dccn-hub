.PHONY: proto build test run deploy clean local

GOPATH:=$(shell go env GOPATH)

all: run test deploy devtools

proto:
	@protoc -I. \
	    -I${GOPATH}/src:. \
	    --micro_out=. \
	    --go_out=. \
	    proto/v1/usermgr.proto

build: proto
	docker build . -t app_dccn_usermgr:v1.0

run: build
	docker run -d \
		-p 50051 \
		$(MICRO_ENV) \
		$(PROGRAM_ENV) \
		app_dccn_usermgr

local: proto
	@$(MICRO_ENV) \
	$(PROGRAM_ENV) \
	go run main.go

clean:
	rm app_dccn_usermgr

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
	MICRO_REGISTRY=mdns \
	MICRO_SERVER_NAME="go.micro.srv.v1.usermgr" \
	MICRO_REGISTRY_ADDRESS="localhost:8500" \
	MICRO_SERVER_VERSION="v1.0" \
	MICRO_REGISTER_TTL=15 \
	MICRO_REGISTER_INTERVAL=15
endef

define PROGRAM_ENV
	DB_HOST="127.0.0.1:27017" \
	DB_NAME="dccn" \
	DB_COLLECTION="user" \
	DB_TIMEOUT=5 \
	DB_POOL_LIMIT=4096
endef
