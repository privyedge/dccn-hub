
GOPATH:=$(shell go env GOPATH)
SUBDIRS := $(wildcard */*.proto)

proto: $(SUBDIRS)
$(SUBDIRS):
	@protoc -I/usr/local/include -I. \
	    -I${GOPATH}/src \
	    -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	    --micro_out=. --go_out=. \
	    $@

.PHONY: proto $(SUBDIRS)

.PHONY: build
build: proto

	go build -o app_dccn_dccenter main.go plugin.go

.PHONY: run
run: build
	MICRO_BROKER=rabbitmq \
	./app_dccn_dccenter

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t app_dccn_dccenter:latest
