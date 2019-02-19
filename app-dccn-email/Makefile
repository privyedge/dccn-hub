
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/v1/email.proto

.PHONY: build
build: proto

	go build -o app_dccn_email main.go plugin.go

.PHONY: local
local:
	MICRO_REGISTRY=consul MICRO_BROKER="rabbitmq" go run main.go


.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t app_dccn_email-srv:latest
