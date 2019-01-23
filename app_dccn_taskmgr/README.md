# Taskmgr Service

This is the Taskmgr service

Generated with

```bash
micro new github.com/Ankr-network/refactor/app_dccn_taskmgr --namespace=go.micro --alias=taskmgr --type=srv
```

## Getting Started

- publisher: publish the v1's task info to "topic.task.new, topic.task.cancel, topic.task.update"
- subscriber: subscribe tasks's result from "topic.task.result"
- handler: request handler
- config: all configuration
- proto: serve's interface

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.v1
- Type: srv
- Alias: v1

## Dependencies

`go get -v go.etcd.io/etcd`
Micro services depend on service discovery. The default is consul.

```bash
# install etcd
brew install etcd

# run etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```bash
make build
```

Run the service

```bash
./taskmgr-srv
```

Build a docker image

```bash
make docker
```
