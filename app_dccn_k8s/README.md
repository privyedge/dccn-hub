# K8s Service

This is the K8s service

Generated with

```
micro new github.com/Ankr-network/refactor/app_dccn_k8s --namespace=network.ankr --alias=k8s --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: network.ankr.srv.k8s
- Type: srv
- Alias: k8s

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
./k8s-srv
```

Build a docker image
```
make docker
```