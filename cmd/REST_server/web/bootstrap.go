package web

import (
	"github.com/gorilla/mux"
)

func Startup(router *mux.Router) {
	usersController := NewUserController()
	router.HandleFunc("/users", usersController.GetAll).Methods("GET")
	router.HandleFunc("/users/{userId}", usersController.GetById).Methods("GET")
	router.HandleFunc("/users", usersController.Create).Methods("POST")
	router.HandleFunc("/users/{userId}", usersController.Update).Methods("PUT")
	router.HandleFunc("/users/{userId}", usersController.Delete).Methods("DELETE")
}
