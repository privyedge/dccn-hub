.PHONY: build
build:
	go build -o taskmgr-cli publisher.go plugin.go

.PHONY: run
run: build
	MICRO_BROKER=rabbitmq \
	./taskmgr-cli

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t taskmgr-srv:latest
