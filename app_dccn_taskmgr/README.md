# Taskmgr Service

This is the Taskmgr service

Generated with

```
micro new github.com/Ankr-network/refactor/app_dccn_taskmgr --namespace=go.micro --alias=taskmgr --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.taskmgr
- Type: srv
- Alias: taskmgr

## Dependencies

Micro services depend on service discovery. The default is consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./taskmgr-srv
```

Build a docker image
```
make docker
```