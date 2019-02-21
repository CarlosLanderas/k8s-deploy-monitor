package main

import (
	"devops-spain/api"
	"devops-spain/k8s"
	"encoding/json"
	"fmt"
	"k8s.io/api/apps/v1"
	"log"
	"net/http"
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

	go deployment.StartWatcher(deploymentChanges)
	go deploymentsHub.Run()

	fmt.Println("Starting API on port", port)

	log.Fatal(http.ListenAndServe(port, srv.Router()))
}
