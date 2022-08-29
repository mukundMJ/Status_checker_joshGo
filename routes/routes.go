package routes

import (
	"github.com/gorilla/mux"
	"josh/site_status_checker/controllers"
)

var RegisterRoutes = func(router *mux.Router) {
	router.HandleFunc("/websites", controllers.AddSites).Methods("POST")
	router.HandleFunc("/websites", controllers.GetSiteStatus).Methods("GET")
}
