package main

import (
	"log"
	"net/http"
	"time"

	"github.com/casmelad/GlobantPOC/cmd/REST_server/web"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	web.Startup(router)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
