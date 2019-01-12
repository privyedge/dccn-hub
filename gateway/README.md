# Api Service

This is the Api service

Generated with

```bash
micro new github.com/Ankr-network/dccn-hub/api --namespace=go.micro --type=api
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.gateway.gateway
- Type: gateway
- Alias: gateway

## Dependencies

Micro services depend on service discovery. The default is consul.

```bash
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```bash
make build
```

Run the service

```bash
./api-api
```

Build a docker image

```bash
make docker
```