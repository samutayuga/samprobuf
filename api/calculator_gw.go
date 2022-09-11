package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/samutayuga/samprobuf/pb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

var (
	clientPort       int
	svrPort          int
	serviceName      string
	assemblyTemplate *template.Template
	configPath       string
)

const (
	connStr = "%s:%d"
)

func main() {
	log.Default().SetFlags(log.LstdFlags | log.Lshortfile)
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	cs := fmt.Sprintf(connStr, serviceName, svrPort)
	log.Printf("dialing server with connection %s", cs)
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
	routes.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods("GET")
	routes.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods("GET")

	fmt.Printf("Calculator Gateway is running on %d\n", clientPort)
	s := fmt.Sprintf(":%d", clientPort)
	http.ListenAndServe(s, routes)

}

type Capability struct {
	CapabilityType string `yaml:"capabilityType"`
}
type Assembly struct {
	SessionId    string       `yaml:"session-id"`
	AssemblyType string       `yaml:"assembly-type"`
	Used         bool         `yaml:"connectedToDs"`
	Capabilities []Capability `yaml:"capabilities"`
}

// InitTemplate ...
func InitTemplate() {
	var err error
	assemblyTemplate, err = template.ParseGlob("*.gotmpl")
	//remember to remove
	if err != nil {
		log.Printf("error while parsing the template %s %v", "*.gotmpl", err)
		if assemblyTemplate, err = template.ParseFiles("assembly.gotmpl"); err != nil {
			log.Fatalf("error while parsing the files %s %v", "*.gotmpl", err)
		}
	}
}
func init() {

	flag.StringVar(&configPath, "configPath", "config/client.yaml", "Provide the path of config file, eg. config/server.yaml")
	flag.Parse()
	InitTemplate()
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile(configPath)
	//v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("problem reading config %v", err)

	}
	appName := v.GetString("application-name")
	connectedStatus := v.Get("connectedToDs")

	if connectedStatus == nil {
		log.Printf("connection status %v\n", false)
	} else {
		connectedStatus = v.GetBool("connectedToDs")
		log.Printf("connection status %v\n", connectedStatus)
	}
	clientPort = v.GetInt("client.port")
	serviceName = v.GetString("server.service-name")
	svrPort = v.GetInt("server.port")
	//read config
	assembly := Assembly{}
	assembly.getAssembly()
	log.Printf("%s get port number from config %d", appName, clientPort)

}
func (a *Assembly) getAssembly() {
	if yamlF, err := os.ReadFile(configPath); err == nil {
		//unmarshal
		if errUnmarshall := yaml.Unmarshal(yamlF, a); errUnmarshall != nil {
			log.Fatalf("error unmarshalling config %v", errUnmarshall)
		} else {
			//display it
			a.display()
		}
	} else {
		log.Fatalf("error reading config %v", err)
	}

}
func (a *Assembly) display() {
	log.Printf("assembly %v\n", a)
	if err := assemblyTemplate.ExecuteTemplate(os.Stdout, "assembly.gotmpl", a); err != nil {
		log.Printf("error while executing template %s for data %v %v", "assembly.gotmpl", a, err)
	}
}
