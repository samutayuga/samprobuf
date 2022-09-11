package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"github.com/samutayuga/samprobuf/pb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	clientPort  int
	svrPort     int
	serviceName string
)

const (
	connStr = "%s:%d"
)

func init() {
	var configPath string
	flag.StringVar(&configPath, "configPath", "config/client.yaml", "Provide the path of config file, eg. config/server.yaml")
	flag.Parse()
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile(configPath)
	//v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("problem reading config %v", err)

	}

	clientPort = v.GetInt("client.port")
	serviceName = v.GetString("server.service-name")
	svrPort = v.GetInt("server.port")

	log.Printf("get port number from config %d", clientPort)

}
func main() {
	log.Default().SetFlags(log.LstdFlags | log.Lshortfile)
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	cs := fmt.Sprintf(connStr, serviceName, svrPort)
	conn, err := grpc.Dial(cs, opts)
	if err != nil {
		log.Fatalf("error while dialing server %d %v", svrPort, err)
	}
	defer conn.Close()
	client := pb.NewPrimerCalculatorClient(conn)
	//ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	//defer cancel()

	cRes, cErr := client.Calculate(context.TODO(), &pb.CalculationRequest{Requestor: "sam", Input: int32(345)})
	if cErr != nil {
		log.Fatalf("Error while calling the calculate method %v", cErr)
	}
	log.Printf("got response %v", cRes)

}
