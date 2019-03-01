package routes

import (
	deployment "k8s-deploy-monitor/k8s"
	"encoding/json"
	"net/http"
)

func GetDeployments(kubeconfig string) func(w http.ResponseWriter, r *http.Request) {

	var deploymentsClient = deployment.NewClient(kubeconfig)

	return func(w http.ResponseWriter, r *http.Request) {
		deployments := deploymentsClient.GetDeployments()
		content, err := json.Marshal(deployments)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(content)
	}
}
