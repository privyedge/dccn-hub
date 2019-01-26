# App_dccn_email Service

This is the App_dccn_email service

Generated with

```
micro new github.com/Ankr-network/dccn-hub/app_dccn_email --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.app_dccn_email
- Type: srv
- Alias: app_dccn_email

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
./app_dccn_email-srv
```

Build a docker image
```
make docker
```