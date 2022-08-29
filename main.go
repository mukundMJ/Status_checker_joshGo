package main

import (
	"github.com/gorilla/mux"
	"josh/site_status_checker/routes"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterRoutes(r)
	http.Handle("/", r)

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}
