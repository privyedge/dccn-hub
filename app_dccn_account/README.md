# Account Service

This is Account Management & Authentication Services

## Contents
- config     - configuration information for hot updates
- db_account - wrapper operations of account table
- handler    - an RPC account service
- proto      - declaration of data structures and interfaces
- testdata   - here are test data
- token      - Authentication Service


## Getting

- [Setup](#setup)
- [Dependencies](#dependencies)
- [Protobufs](#protobufs)
- [Warnning](#warnning)
- [TODO](#todo)


Generated with

```
micro new github.com/Ankr-network/refactor/app_dccn_account --namespace=network.ankr --alias=account --type=srv
```

### Setup

Docker is required for running the services https://docs.docker.com/engine/installation.

Protobuf v3 are required:

    $ brew install protobuf

Install the protoc-gen libraries and other dependencies:

    $ go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
    $ go get -u github.com/micro/protoc-gen-micro
    $ go get -u github.com/micro/go-micro
    $ go get -u github.com/hailocab/go-geoindex

## Dependencies

    $ go get -u -d github.com:micro/micro
    $ go get -u -d gopkg.in/mgo.v2
    $ go get -u -d github.com/dgrijalva/jwt-go

### Protobufs

  If changes are made to the Protocol Buffer files use the Makefile to regenerate:

      $ make proto

  ### Run

  To make the demo as straigforward as possible; [Docker Compose](https://docs.docker.com/compose/) is used to run all the services at once (In   a production environment each of the services would be run (and scaled) independently).

      $ make build
      $ make run


### Warnning
    $ Account in database can only be operated and maintained by the service

### TODO:
    $ using TOML as a configuration file
    $ add support for bulk operations
    $ support for filtering lookups
