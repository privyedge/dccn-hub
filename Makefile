.PHONY: build test run deploy local clean

GOPATH:=$(shell go env GOPATH)

all: run test deploy devtools

#=================== BUILD ===================
build-api:
	-docker build -f dockerfile/api.Dockerfile -t dccn-api:latest .

build-task:
	-docker build -f dockerfile/task.Dockerfile -t dccn-task:latest .

build-user:
	-docker build -f dockerfile/user.Dockerfile -t dccn-user:latest .

build-email:
	-docker build -f dockerfile/email.Dockerfile -t dccn-email:latest .

build-test:
	-docker build -f dockerfile/test.Dockerfile -t dccn-test-all:latest .

build: build-api build-email build-task build-task build-test

#=================== CREATE ===================
create-dep:
	-kubectl create -f kubernetes/broker.yml
	-kubectl create -f kubernetes/consul.yml
	-kubectl create -f kubernetes/datastore.yml
	kubectl get pod

create-app:
	-kubectl create -f kubernetes/task.yml
	-kubectl create -f kubernetes/user.yml
	-kubectl create -f kubernetes/email.yml
	-kubectl create -f kubernetes/api.yml

create: build create-dep create-app
	kubectl get pod

#=================== RUN ===================
run: build
	docker run -d \
		-p 50051 \
		app_dccn_dcmgr

#=================== TEST ===================
test:
	kubectl create -f kubernetes/test.yml

#=================== CLEAN UP ===================
clean-api:
	-kubectl delete -f kubernetes/api.yml

clean-task:
	-kubectl delete -f kubernetes/task.yml

clean-user:
	-kubectl delete -f kubernetes/user.yml

clean-email:
	-kubectl delete -f kubernetes/email.yml

clean-test:
	-kubectl delete -f kubernetes/test.yml

clean-broker:
	-kubectl delete -f kubernetes/broker.yml

clean-datastore:
	-kubectl delete -f kubernetes/datastore.yml

clean-consul:
	-kubectl delete -f kubernetes/consul.yml

clean-dep: clean-broker clean-datastore clean-consul
	kubectl get pod

clean-app: clean-test clean-api clean-broker clean-email clean-task clean-user
	kubectl get pod

clean: clean-dep clean-app
	kubectl get pod

#=================== DEPLOY ===================
deploy:
	@echo "docker push"

#=================== DEPENDENCY ===================
devtools:
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go get github.com/micro/protoc-gen-micro
	@type "protoc" 2> /dev/null || echo 'Please install protoc'
	@type "protoc-gen-micro" 2> /dev/null || echo 'Please install protoc-gen-micro'

