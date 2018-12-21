
GOPATH := $(shell go env GOPATH)
#SUBDIRS := $(wildcard ../app_*/proto/*.proto)
SUBDIRS := $(wildcard */*/*.proto)

.PHONY: proto
proto: $(SUBDIRS)
$(SUBDIRS):
	@protoc -I/usr/local/include -I. \
	  	-I${GOPATH}/src \
	  	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  	--go_out=plugins=grpc:. \
	  	$@
	@protoc -I/usr/local/include -I. \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		$@
	@protoc -I. \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=. \
		$@

.PHONY: proto $(SUBDIRS)

.PHONY: build
build: proto

	go build -o gateway main.go

.PHONY: run
run: build
	./gateway

.PHONY: clean
clean:
	rm gateway

# .PHONY: test
# test:
# 	go test -v ./... -cover

# .PHONY: docker
# docker:
# 	docker build . -t taskmgr-srv:latest
