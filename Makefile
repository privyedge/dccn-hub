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
	-kubectl create -f deployments/broker.yml
	-kubectl create -f deployments/consul.yml
	-kubectl create -f deployments/datastore.yml
	kubectl get pod

create-app:
	-kubectl create -f deployments/task.yml
	-kubectl create -f deployments/user.yml
	-kubectl create -f deployments/email.yml
	-kubectl create -f deployments/api.yml

create-test:
	-kubectl create -f deployments/test.yml

create: build create-dep create-app
	kubectl get pod

#=================== RUN ===================
run: build
	docker run -d \
		-p 50051 \
		app_dccn_dcmgr

#=================== TEST ===================
test:
	kubectl create -f deployments/test.yml

#=================== CLEAN UP ===================
clean-api:
	-kubectl delete deployment.apps/api

clean-task:
	-kubectl delete deployment.apps/api

clean-user:
	-kubectl delete -f deployments/user.yml

clean-email:
	-kubectl delete -f deployments/email.yml

clean-test:
	-kubectl delete -f deployments/test.yml

clean-broker:
	-kubectl delete -f deployments/broker.yml

clean-datastore:
	-kubectl delete -f deployments/datastore.yml

clean-consul:
	-kubectl delete -f deployments/consul.yml

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

