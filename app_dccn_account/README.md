# Dcmgr Service

This is the Dcmgr service

Generated with

```
micro new github.com/Ankr-network/refactor/app_dccn_dcmgr --namespace=network.ankr --alias=dcmgr --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: network.ankr.srv.dcmgr
- Type: srv
- Alias: dcmgr

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
./dcmgr-srv
```

Build a docker image
```
make docker
```