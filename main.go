package main

import (
	"encoding/json"
	"fmt"
	"k8s-deploy-monitor/api"
	"k8s-deploy-monitor/k8s"
	"log"
	"net/http"

	"k8s.io/api/apps/v1"
)

const port = ":8080"

func main() {

	deploymentChanges := make(chan *v1.Deployment)
	srv := server.Create()

	deploymentsHub := srv.DeploymentsHub()
	go func() {
		for message := range deploymentChanges {
			jsonMessage, _ := json.Marshal(message)
			deploymentsHub.Broadcast <- []byte(jsonMessage)
		}
	}()

	deploymentsClient, err := deployment.NewFromConfig()

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			deploymentsClient.StartWatcher(deploymentChanges)
		}
	}()

	go deploymentsHub.Run()

	fmt.Println("Starting API on port", port)

	log.Fatal(http.ListenAndServe(port, srv.Router()))
}
