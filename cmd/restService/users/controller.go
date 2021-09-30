package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usersController struct {
	dataSource *UserProxy
}

func (u *usersController) GetAll(w http.ResponseWriter, r *http.Request) {

	resp, err := u.dataSource.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := resp

	respondWithJSON(w, http.StatusOK, response)

}

func (u *usersController) GetById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	email, ok := vars["email"]

	if !ok {
		respondWithError(w, http.StatusBadRequest, "invalida input data")
		return
	}

	user, err := u.dataSource.GetByEmail(email)

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (u *usersController) Create(w http.ResponseWriter, r *http.Request) {

	userToCreate := User{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&userToCreate); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	userCreated, err := u.dataSource.Create(userToCreate)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, userCreated)
}

func (u *usersController) CreateMany(w http.ResponseWriter, r *http.Request) {

	for i := 0; i < 100; i++ {

		go func(index int) {
			intId := strconv.Itoa(index)

			_, err := u.dataSource.Create(User{
				Email:    "user" + intId + "@gmail.com",
				Name:     "test",
				LastName: "test_last",
			})

			if err != nil {
				fmt.Println(err.Error())
			}

		}(i)
	}

	respondWithJSON(w, http.StatusNoContent, "Ok")
}

func (u *usersController) Update(w http.ResponseWriter, r *http.Request) {

	userToUpdate := User{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&userToUpdate); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	vars := mux.Vars(r)

	email, ok := vars["email"]

	if !ok {
		respondWithError(w, http.StatusBadRequest, "invalida input data")
		return
	}

	userToUpdate.Email = email

	defer r.Body.Close()

	if _, err := u.dataSource.Update(userToUpdate); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, userToUpdate)
}

func (u *usersController) Delete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, ok := vars["userId"]

	if !ok {
		respondWithError(w, http.StatusBadRequest, "invalida input data")
		return
	}

	intId, res := strconv.Atoi(id)

	if res != nil {
		respondWithError(w, http.StatusBadRequest, "invalida input data")
		return
	}

	user, err := u.dataSource.Delete(int(intId))

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, user)
}

func NewUserController() *usersController {

	return &usersController{
		dataSource: NewUserProxy(),
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
