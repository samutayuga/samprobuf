package main

import (
	"context"

	"log"
	"net"

	"github.com/samutayuga/samprobuf/calculator/primerpb"
	grpc "google.golang.org/grpc"
)

type Calcserver struct {
}

func (s *Calcserver) Calculate(ctx context.Context, in *primerpb.CalculationRequest) (resp *primerpb.CalculationResponse, err error) {
	r := &primerpb.CalculationResponse{Message: "test", Count: int32(10)}

	return r, nil
}

func main() {

	log.Println("calculator server is started..")
	l, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("error while creating listener %v", err)
	}
	//register the service
	s := grpc.NewServer()
	primerpb.RegisterPrimerCalculatorServer(s, &Calcserver{})
	

}
