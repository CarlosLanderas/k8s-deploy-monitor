package main

import (
	server "devops-spain/api"
	deployment "devops-spain/k8s"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	ini "gopkg.in/ini.v1"
	v1 "k8s.io/api/apps/v1"
)

const port = ":8080"

func main() {

	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("Could not read config.ini")
	}

	kubeconfig := cfg.Section("").Key("kubeconfig").Value()

	deploymentChanges := make(chan *v1.Deployment)
	srv := server.Create(kubeconfig)

	deploymentsHub := srv.DeploymentsHub()
	go func() {
		for message := range deploymentChanges {
			jsonMessage, _ := json.Marshal(message)
			deploymentsHub.Broadcast <- []byte(jsonMessage)
		}
	}()

	deploymentsClient := deployment.NewClient(kubeconfig)

	go func() {
		for {
			deploymentsClient.StartWatcher(deploymentChanges)
		}
	}()

	go deploymentsHub.Run()

	fmt.Println("Starting API on port", port)

	log.Fatal(http.ListenAndServe(port, srv.Router()))
}
