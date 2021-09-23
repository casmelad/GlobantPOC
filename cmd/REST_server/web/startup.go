package web

import (
	"github.com/gorilla/mux"
)

func Startup(router *mux.Router) {
	usersController := NewUserController()
	router.HandleFunc("/users", usersController.GetAll).Methods("GET")
	router.HandleFunc("/", usersController.Hello).Methods("GET")
	router.HandleFunc("/users/{email}", usersController.GetById).Methods("GET")
	router.HandleFunc("/users", usersController.Create).Methods("POST")
	router.HandleFunc("/users/{email}", usersController.Update).Methods("PUT")
	router.HandleFunc("/users/{userId}", usersController.Delete).Methods("DELETE")
	router.HandleFunc("/users/createmany", usersController.CreateMany).Methods("POST")
}
