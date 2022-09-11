# Micro services with go application
> As a user I want to have a service which provide an API to calculate the maximum number of prime number for a given number

`service definition`

The communication protocol used is `gRPC`.
The following is the defintion of service and messages in  `protocol buffer`.

```proto
// The calculator service defintion
// a unary API
// Synchronous API, request-response model

syntax = "proto3";
package pb;
option go_package="calc/pb";
// The calculator service defintion
// a unary API
// Synchronous API, request-response model
service PrimerCalculator{
    rpc Calculate(CalculationRequest) returns (CalculationResponse){}
}

message CalculationRequest {
string requestor=1;
    int32 input=2;
}

message CalculationResponse {
string message=1;
}
```
Compile them using a command as below,

```shell
protoc -I . \
   --go_out ./pb --go_opt paths=source_relative \
   --go-grpc_out ./pb --go-grpc_opt paths=source_relative \
   ./*.proto
```
The `proto compiler` will generate 2 files,
* <proto_file>_grpc.go
  > the grpc services stubs
* <proto_file>_pb.go
  > the messages go codes


## Code Base
This is an API to provide the calculation of prime number
* Generated go codes
  > Do not touch this files
* Service Implementation
  > The main changes will be in this part. Need a `main` function to launch an http server. 
  > Implement the stubs. Basically, implement the interface `PrimerCalculatorServer` for the defined method,
  ```go
  type PrimerCalculatorServer interface {
	Calculate(context.Context, *CalculationRequest) (*CalculationResponse, error)

	mustEmbedUnimplementedPrimerCalculatorServer()
  }
  ```
  This is the server microservices


* Service Gateway
  > Implement an http server to receive the `RESTfull` request from external then route to the server though `gRPC` protocol
  This the gateway microservices

## Containerization
### Service
`Dockerfile`
```docker
# syntax=docker/dockerfile:1
FROM golang:1.18-buster AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download


COPY ./service/main.go ./
COPY ./pb ./pb
RUN go build -o /calculator
# Final Stage - Stage 2
FROM gcr.io/distroless/base-debian10 as baseImage
WORKDIR /app
COPY --from=builder /calculator ./calculator
COPY ./service/config ./config

ENTRYPOINT ["/app/calculator"]
EXPOSE 8001
```
Building

```shell
docker build -t samutup/calculator:1.0.0 --no-cache -f DockerFile-calculator .
```

### Gateway
`Dockerfile`
```docker
# syntax=docker/dockerfile:1
FROM golang:1.18-buster AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download


COPY ./api/calculator_gw.go ./
COPY ./pb ./pb
RUN go build -o /calculator-gw
# Final Stage - Stage 2
FROM gcr.io/distroless/base-debian10 as baseImage
WORKDIR /app
COPY --from=builder /calculator-gw ./calculator-gw
COPY ./api/config ./config
COPY ./api/assembly.gotmpl ./
ENTRYPOINT ["/app/calculator-gw"]
EXPOSE 8002
```
Building

```shell
docker build -t samutup/calculator_gw:1.0.0 --no-cache -f DockerFile-calculator-gw .
```

## Kubernetes Objects Defnition File

### Service
`ConfigMaps`
> To manage the configuration used by the `container` or `kubernetes` object

`Deployment`

`Pod`
```yaml

```

`Service`
> Exposing the service, in term of `hostname:port` to external. Can be outside or inside the cluster


`ReplicaSet`

`Ingress`


## Helm packaging
> Package manager
`helm chart`


