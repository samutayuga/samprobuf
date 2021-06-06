# Micro services with go application
> As a user I want to have a service which provide an API to calculate the maximum number of prime number for a given number
`service definition`
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

## Code Base
This is an API to provide the calculation of prime number

`Proto File`

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/primer_calculator.proto
```

## Containerization

`Dockerfile`

## Kubernetes Objects Defnition File
## Helm packaging
`helm chart`


