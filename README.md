# dccn-hub
This is the central component for Ankr DCCN. Ankr Hub consists of two microservices: k8s_handler and task_manager.

## k8s_handler

* Receive DC status and update to database
* deliver task to DC when
    1.  receive new task notification from task_manager service  
    2. interval query task for new task
* receive task status from DC and update to database

## task_manager

* add task from cli to database
* query task list from database and return to cli
* notify k8s_handler when new task create

### Some tech thoughtsß
* Database
    - use MongoDB single instance
* Collections
    - task
    - user
    - datacenter
* Ankr Hub connect with cli/k8s by gRPC (k8s may use ZeroMQ as messaging)

## Workflow

* Users define desired infrastructure and workloads to run on the remote infrastructure, and how workloads can connect to one another.
* Desired lifetime of resources is expressed via collateral requirements.
* Orders are generated from the deployment requirement.
* Datacenters bid on an open orderbook.
* The bid with lowest price (more heuristics later) gets matched with order to create a contract.
* Once a contract is reached, workloads and topology are delivered to the datacenter(s).
* Datacenter(s) deploy workloads and allow connectivity as specified by the tenant.
* If a datacenter fails to maintain contract, collateral is transferred to tenant, and a new order is created for the desired resources.
* A tenant can close any active deployment at any time

## Install
* set $GOPATH

* install grcp package  
  * go get -u google.golang.org/grpc

* git clone code  
  * cd $GOPATH/src  
  * git clone git@github.com:Ankr-network/dccn-hub.git  -b feature/78-ankr-hub dccn-hub

* run server :   
  * go run taskmanager/service.go

* run client:   
  * go run test/cli/add_task.go

* install MongoDB  
  * https://treehouse.github.io/installation-guides/mac/mongo-mac.html

* New Way Run MongoDB  (by docker)   
  * cd dccn-hub/docker/   
  * docker run   -p 27017:27017  --name ankr_mongo -d mongo  
  * docker logs ankr_mongo  // check logs

* Two way to install default data:
  1. go run db/install.go  
  2. mongorestore -d test db/backup   


* To test mongo is running
  * mongo   
  * use test    
  * db.user.find()

* proto compiler tools
  * go get github.com/golang/protobuf/protoc-gen-go   
  * protoc --go_out=plugins=grpc:. *.proto

## Building with Docker and CircleCI
Using the docker build using the "Dockerfile.dep" file if you download the source and build locally:
```
dep ensure -update
docker build -f local.dockerfile. -t dccn_hub .
docker run -p 50051:50051 dccn_hub
```

for the CircleCI setting, check the .circleci/config.yml for detail,  CircleCI pipeline will build and push the docker image to aws ecr repository "815280425737.dkr.ecr.us-west-2.amazonaws.com/dccn_ecr"
