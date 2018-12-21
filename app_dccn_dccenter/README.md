# Dccenter Service

This is the Dccenter service

Generated with

```
micro new github.com/Ankr-network/refactor/app_dccn_dccenter --namespace=network.ankr --alias=dccenter --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: network.ankr.srv.dccenter
- Type: srv
- Alias: dccenter

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
./dccenter-srv
```

Build a docker image
```
make docker
```