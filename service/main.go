package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/samutayuga/samprobuf/pb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	port    int
	svrType string
)

const (
	msgNotPrime = "The number %d requested by %s is not a prime"
	msgPrime    = "The number %d requested by %s is a prime"
)

func IsPrime(aNumber int) bool {
	if aNumber == 2 {
		return true
	}
	//if even
	if aNumber > 2 && aNumber%2 == 0 {
		return false
	}
	//if odd
	for i := 3; float64(i) <= math.Sqrt(float64(aNumber)); i += 2 {
		if aNumber%i == 0 {
			log.Printf("Found at least one divider other than %d which is %d", aNumber, aNumber/i)
			return false
		}
	}
	return true
}

type Calcserver struct {
	pb.UnimplementedPrimerCalculatorServer
}

func (s *Calcserver) Calculate(ctx context.Context, in *pb.CalculationRequest) (resp *pb.CalculationResponse, err error) {
	log.Printf("serving one request from %s\n", in.Requestor)

	if IsPrime(int(in.Input)) {
		s := fmt.Sprintf(msgPrime, in.Input, in.Requestor)
		r := &pb.CalculationResponse{Message: s}
		return r, nil
	} else {
		s := fmt.Sprintf(msgNotPrime, in.Input, in.Requestor)
		r := &pb.CalculationResponse{Message: s}
		return r, nil
	}
}
func main() {

	log.Default().SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("calculator server is started..")
	svrStr := fmt.Sprintf(":%d", port)
	l, err := net.Listen(svrType, svrStr)
	if err != nil {
		log.Fatalf("error while creating listener %v", err)
	}
	//register the service
	s := grpc.NewServer()
	pb.RegisterPrimerCalculatorServer(s, &Calcserver{})
	log.Println("Service is registered...")
	go func() {
		fmt.Println("Server starting ....")
		if err := s.Serve(l); err != nil {
			log.Fatalf("failed to server %v", err)
		}

	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGKILL)
	<-ch
	fmt.Println("Stopping the server...")
	s.Stop()
	fmt.Println("Closing the listener")
	l.Close()
}
func init() {
	var serverConfigPath string
	flag.StringVar(&serverConfigPath, "serverConfig", "config/server.yaml", "Provide the path of config file, eg. config/server.yaml")
	flag.Parse()
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile(serverConfigPath)
	//v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("problem reading config %v", err)

	}

	port = v.GetInt("server.port")
	svrType = v.GetString("server.type")
	log.Printf("get port number from config %d", port)

}
