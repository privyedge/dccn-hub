# Micro gRPC [![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0) [![GoDoc](https://godoc.org/github.com/micro/go-grpc?status.svg)](https://godoc.org/github.com/micro/go-grpc) [![Travis CI](https://api.travis-ci.org/micro/go-grpc.svg?branch=master)](https://travis-ci.org/micro/go-grpc) [![Go Report Card](https://goreportcard.com/badge/micro/go-grpc)](https://goreportcard.com/report/github.com/micro/go-grpc)

A micro gRPC framework. A simplified experience for building gRPC services.

## Overview

Go gRPC makes use of [go-micro](https://github.com/micro/go-micro) plugins to create a simpler framework for gRPC development. It interoperates with
standard gRPC services seamlessly, including the [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway). The go-grpc library uses
the go-micro client and server plugins which make use of [github.com/grpc/grpc-go](https://github.com/grpc/grpc-go). This means we ignore
the go-micro codec and transport but you get a native grpc experience.

<img src="https://micro.mu/docs/images/go-grpc.png" />

## Examples

Find an example greeter service in [examples/greeter](https://github.com/micro/go-grpc/tree/master/examples/greeter).

## Getting Started

See the [docs](https://micro.mu/docs/go-grpc.html) to get started.
