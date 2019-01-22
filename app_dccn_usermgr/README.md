# user Service

This is v1 Management & Authentication Services

## Contents

- config     - configuration information for hot updates
- db_service - wrapper operations of v1 table
- handler    - an RPC v1 service
- proto      - declaration of data structures and interfaces
- testdata   - here are test data
- token      - Authentication Service

## Getting

- [Setup](#setup)
- [Dependencies](#dependencies)
- [Protobufs](#protobufs)
- [Warnning](#warnning)
- [TODO](#todo)
- [Curl](#curl)

Generated with

```bash
micro new github.com/Ankr-network/refactor/app_dccn_user --namespace=network.ankr --alias=user --type=srv
```

### Setup

Docker is required for running the services `https://docs.docker.com/engine/installation`.

Protobuf v3 are required:

```bash
brew install protobuf
```

Install the protoc libraries and other dependencies:

```bash
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u github.com/micro/protoc-gen-micro
go get -u github.com/micro/go-micro
go get -u github.com/hailocab/go-geoindex
```

## Dependencies

    go get -u -d gopkg.in/mgo.v2
    go get -u -d github.com/dgrijalva/jwt-go
    go get -u -d github.com/micro/micro
    go get -u -d github.com/micro/go-micro
    go get github.com/micro/kubernetes
    go get github.com/micro/go-plugins

### Protobufs

  If changes are made to the Protocol Buffer files use the Makefile to regenerate:

```bash
make proto
```

### Run

  To make the demo as straigforward as possible; [Docker Compose](https://docs.docker.com/compose/) is used to run all the services at once (In   a production environment each of the services would be run (and scaled) independently).
      $ make build
      $ make run

### Warnning

- user in database can only be operated and maintained by the service

### TODO

- using TOML as a configuration file
- add support for bulk operations
- support for filtering lookups

### Curl

- Authentication request

```bash
curl --request POST http://host:port/session -H "Content-Type:application/json" -d '{"username":"admin","password":"123"}'
```

- Resourse request

```bash
curl --request GET http://www.video4.cn:3030/api/projects/2 -H "Authorization:Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6MSwidXNlcm5hbWUiOiJhZG1pbiJ9.d1bf66192c1bff9038bcd212ba05dfde55c40d4e2254dd99c9c7653dd27c39ba"
```