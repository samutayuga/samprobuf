package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/samutayuga/samprobuf/calculator/primerpb"
	"github.com/samutayuga/samprobuf/service/decomposer"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	port int
)

const (
	msgNotPrime = "The number %d requested by %s is not a prime"
	msgPrime    = "The number %d requested by %s is a prime"
)

type Calcserver struct {
	primerpb.UnimplementedPrimerCalculatorServer
}

func (s *Calcserver) Calculate(ctx context.Context, in *primerpb.CalculationRequest) (resp *primerpb.CalculationResponse, err error) {
	log.Printf("serving one request from %s\n", in.Requestor)

	if decomposer.IsPrime(int(in.Input)) {
		s := fmt.Sprintf(msgPrime, in.Input, in.Requestor)
		r := &primerpb.CalculationResponse{Message: s, Count: in.Input}
		return r, nil
	}
	if !decomposer.IsPrime(int(in.Input)) {
		s := fmt.Sprintf(msgNotPrime, in.Input, in.Requestor)
		r := &primerpb.CalculationResponse{Message: s, Count: in.Input}
		return r, nil
	}
	return nil, nil
}
func main() {
	log.Default().SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("calculator server is started..")
	svrStr := fmt.Sprintf("0.0.0.0:%d", port)
	l, err := net.Listen("tcp", svrStr)
	if err != nil {
		log.Fatalf("error while creating listener %v", err)
	}
	//register the service
	s := grpc.NewServer()
	primerpb.RegisterPrimerCalculatorServer(s, &Calcserver{})
	log.Println("Service is registered...")
	go func() {
		fmt.Println("Server starting ....")
		if err := s.Serve(l); err != nil {
			log.Fatalf("failed to server %v", err)
		}

	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println("Stopping the server...")
	s.Stop()
	fmt.Println("Closing the listener")
	l.Close()
}
func init() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile("config/server.yaml")
	//v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("problem reading config %v", err)

	}

	port = v.GetInt("server.port")
	log.Printf("get port number from config %d", port)

}
