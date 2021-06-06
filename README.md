# Micro services with go application
> As a user I want to have a service which provide an API to calculate the maximum number of prime number for a given number
`service definition`

The communication protocol used is `gRPC`.
The foollowing is the defintion of service and messages in  `protcol buffer`.

```proto
// The calculator service defintion
// a unary API
// Synchronous API, request-response model

service PrimerCalculator{
    rpc Calculate(CalculationRequest) returns (CalculationResponse){}
}

message CalculationRequest {
    required string requestor=1;
    required int32 input=2;
}

message CalculationResponse {
    required string message=1;
    required int32 count=2;
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
FROM golang
COPY . /go/src/prime
WORKDIR /go/src/prime
RUN go mod init
RUN go get ./...
ENTRYPOINT go run main.go
EXPOSE 8001
```
Building

```shell
cd service

docker build -t samutup/primer-checker:1.0.0 --no-cache -f DockerFile .
```

### Gateway
`Dockerfile`
```docker
FROM golang
COPY . /go/src/prime
WORKDIR /go/src/prime
RUN go mod init
RUN go get ./...
ENTRYPOINT go run calculator_gw.go
EXPOSE 8002
```
Building

```shell
cd api

docker build -t samutup/primer-gw:1.0.0 --no-cache -f DockerFile .
```

## Kubernetes Objects Defnition File

### Service
```yaml

```

## Helm packaging
`helm chart`


