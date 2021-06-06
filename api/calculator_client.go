package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

func main() {
	log.Default().SetFlags(log.LstdFlags | log.Lshortfile)
	opts := grpc.WithInsecure()
	cs := fmt.Sprintf(connStr, serviceName, svrPort)
	conn, err := grpc.Dial(cs, opts)
	if err != nil {
		log.Fatalf("error while dialing server 8001 %v", err)
	}
	defer conn.Close()
	client := pb.NewPrimerCalculatorClient(conn)
	routes := mux.NewRouter()
	routes.HandleFunc("/prime/{name}/{number}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UFT-8")
		vars := mux.Vars(r)
		if requestor, exists := vars["name"]; !exists {
			log.Println("error while retrieving the name")
			w.WriteHeader(http.StatusBadRequest)
			m := "name is not provided"
			w.Write([]byte(m))
		} else {
			if number, err := strconv.ParseUint(vars["number"], 10, 64); err != nil {
				log.Printf("error while retrieving number %v\n", err)
				m := "number is not provided"
				w.Write([]byte(m))
			} else {
				log.Printf("Get the request from %s to compute %d", requestor, number)
				//call backend
				ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
				defer cancel()

				cRes, cErr := client.Calculate(ctx, &pb.CalculationRequest{Requestor: requestor, Input: int32(number)})
				if cErr != nil {
					log.Fatalf("Error while calling the calculate method %v", cErr)
				}
				log.Printf("got response %v", cRes)

				if b, mErr := json.Marshal(cRes); mErr != nil {
					log.Printf("error while marshalling response to json %v", mErr)
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write(b)
				}

			}
		}

	}).Methods("GET")
	fmt.Printf("Client is running on %d\n", clientPort)
	s := fmt.Sprintf(":%d", clientPort)
	http.ListenAndServe(s, routes)

}
func init() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile("config/client.yaml")
	//v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("problem reading config %v", err)

	}

	clientPort = v.GetInt("client.port")
	serviceName = v.GetString("server.service-name")
	svrPort = v.GetInt("server.port")

	log.Printf("get port number from config %d", clientPort)

}
