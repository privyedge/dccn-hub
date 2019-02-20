
GOPATH:=$(shell go env GOPATH)

.PHONY: build test docker local

build: proto
	go build -o app_dccn_email main.go plugin.go

test:
	go test -v ./... -cover

docker:
	docker build . -t app_dccn_email-srv:latest

local:
	@$(MICRO_ENV) \
	$(PROGRAM_ENV) \
	go run main.go

define MICRO_ENV
	MICRO_REGISTRY=consul \
	MICRO_BROKER=rabbitmq \
	MICRO_SERVER_VERSION="v1.0" \
	MICRO_REGISTER_INTERVAL=15
endef

define PROGRAM_ENV
	DB_HOST="127.0.0.1:27017" \
	DB_NAME="dccn" \
	DB_COLLECTION="user" \
	DB_TIMEOUT=5 \
	DB_POOL_LIMIT=4096 \
	TOKEN_ACTIVE_TIME=20
endef