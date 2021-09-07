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

	usersController := web.NewUserController()

	router.HandleFunc("/users", usersController.GetAll).Methods("GET")
	router.HandleFunc("/users/{userId}", usersController.GetById).Methods("GET")
	router.HandleFunc("/users", usersController.Create).Methods("POST")
	router.HandleFunc("/users/{userId}", usersController.Update).Methods("PUT")
	router.HandleFunc("/users/{userId}", usersController.Delete).Methods("DELETE")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
