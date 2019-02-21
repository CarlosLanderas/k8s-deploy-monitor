package server

import (
	"devops-spain/api/routes"
	"devops-spain/hub"
	"github.com/gorilla/mux"
	"net/http"
)

type Api struct {
	router http.Handler
	deploymentsHub *hub.Hub
}

func (a *Api) Router() http.Handler {
	return a.router
}

func (a *Api) DeploymentsHub() *hub.Hub {
	return a.deploymentsHub
}

//Server interface
type Server interface {
	Router() http.Handler
	DeploymentsHub() *hub.Hub
}

// Create a new api Server
func Create() Server {

	router := mux.NewRouter()

	router.HandleFunc("/deployments", routes.GetDeployments)
	router.Handle("/", http.FileServer(http.Dir("client/build")))

	router.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r,"client/build/manifest.json")
	})

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./client/build/static"))))

	deploymentsHub := hub.NewDeploymentsHub()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(deploymentsHub, w, r)
	})

	return &Api{
		router: router,
		deploymentsHub: deploymentsHub,
	}
}
