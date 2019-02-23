.PHONY: build test run deploy clean local

GOPATH:=$(shell go env GOPATH)

all: run test deploy devtools

build:
	docker build . -t app_dccn_usermgr:v1.0

run: build
	docker run -d \
		-p 50051 \
		$(MICRO_ENV) \
		$(PROGRAM_ENV) \
		app_dccn_usermgr

local:
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
	MICRO_REGISTRY=consul \
	MICRO_BROKER=rabbitmq \
	MICRO_SERVER_VERSION="v1.0" \
	MICRO_REGISTER_INTERVAL=20 \
	MICRO_REGISTER_TTL=30
endef

define PROGRAM_ENV
	DB_HOST="127.0.0.1:27017" \
	DB_NAME="dccn" \
	DB_COLLECTION="user" \
	DB_TIMEOUT=5 \
	DB_POOL_LIMIT=4096 \
	TOKEN_ACTIVE_TIME=20
endef
