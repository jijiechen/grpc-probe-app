
An example application demonstrating how to use gRPC health checking.

# Usage

You may use this app directly using docker:

```
docker run --rm -it jijiechen/grpc-probe-app:latest
```

Or if you'd like to build yourself, please clone the code and execute the following commands:

```sh
 go build -o grpc-probe-app .
 ./grpc-probe-app [--port 5085]
```

## Test the greeter server

Execute these commands:

```sh
# install grpcurl
# go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

grpcurl -proto ./proto/helloworld.proto -d '{"name": "Avatar"}' -plaintext localhost:5085 helloworld.Greeter/SayHello
```

## Test the health check probe server

Execute these commands:

```sh
# install grpcurl
# go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

grpcurl -proto ./proto/healthcheck.proto -d '{"service": "readiness"}' -plaintext localhost:5085 grpc.health.v1.Health/Check
grpcurl -proto ./proto/healthcheck.proto -d '{"service": "liveness"}' -plaintext localhost:5085 grpc.health.v1.Health/Watch
```


# Protocol

1. This application implements the [hello world greeter protocol](https://github.com/grpc/grpc-go/blob/master/examples/helloworld/helloworld/helloworld.proto).

2. This application implements the gRPC health checking protocol. The protocol is defined in the following document:

https://github.com/grpc/grpc/blob/master/doc/health-checking.md

It can be used as a Kubernetes gRPC readiness/liveness/startup probe handler.

# TroubleShooting

Q: Readiness probe errored: missing probe handler for pod(uid):container

A: Please make sure your Kubernetes cluster support the gRPC container probe feature. The feature is available in Kubernetes 1.23+, and requires the `GRPCContainerProbe` feature gate to be turned on a Kubernetes cluster older than 1.27.