package routes

import (
	"devops-spain/k8s"
	"encoding/json"
	"net/http"
)
var deploymentsClient = deployment.NewClient()

func GetDeployments(w http.ResponseWriter, r *http.Request) {
	deployments := deploymentsClient.GetDeployments()
	content, err := json.Marshal(deployments)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(content)
}