package web

import (
	"encoding/json"
	"net/http"

	users "github.com/casmelad/GlobantPOC/cmd/REST_server/application"
	entities "github.com/casmelad/GlobantPOC/cmd/REST_server/entities"
)

type usersController struct {
	dataSource *users.UserProxy
}

func (u *usersController) GetAll(w http.ResponseWriter, r *http.Request) {

	response := []entities.User{}

	resp, err := u.dataSource.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response = resp

	respondWithJSON(w, http.StatusOK, response)

}

func (u *usersController) GetById(r http.ResponseWriter, w *http.Request) {
	user := entities.User{}
	user.Name = "Adrian"
	r.Write([]byte("Not implemented"))
}

func (u *usersController) Create(w http.ResponseWriter, r *http.Request) {

	userToCreate := entities.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userToCreate); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if _, err := u.dataSource.Create(userToCreate); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, userToCreate)
}

func (u *usersController) Update(r http.ResponseWriter, w *http.Request) {
	r.Write([]byte("Not implemented"))
}

func (u *usersController) Delete(r http.ResponseWriter, w *http.Request) {
	r.Write([]byte("Not implemented"))
}

func NewUserController() *usersController {

	return &usersController{
		dataSource: users.NewUserProxy(),
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
