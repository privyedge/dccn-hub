# App_dccn_test Service

This is the App_dccn_test service

Generated with

```
micro new github.com/Ankr-network/dccn-hub/app_dccn_test --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.app_dccn_test
- Type: srv
- Alias: app_dccn_test

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
./app_dccn_test-srv
```

Build a docker image
```
make docker
```