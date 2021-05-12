package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mwazovzky/loki/models"
	"gorm.io/gorm"
)

type UserHandlers struct {
	users *models.Users
}

func NewUserHandlers(db *gorm.DB) *UserHandlers {
	users := models.NewUsers(db)
	return &UserHandlers{users}
}

func (uh *UserHandlers) Index(rw http.ResponseWriter, r *http.Request) {
	log.Println("User Index Request")

	var users []models.User
	uh.users.Find(&users)

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err := e.Encode(users)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (uh *UserHandlers) Show(rw http.ResponseWriter, r *http.Request) {
	log.Println("Users Show Request")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	var user models.User
	uh.users.First(&user, id)

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err = e.Encode(user)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

// curl -X POST localhost:3000/api/users -d '{"name":"test","email":"test@example.com"}'
func (uh *UserHandlers) Create(rw http.ResponseWriter, r *http.Request) {
	log.Println("Users Create Request")

	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	uh.users.Create(&user)

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err = e.Encode(user)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

// curl -X PUT localhost:3000/api/users/5/update -d '{"email":"updated@example.com"}'
func (uh *UserHandlers) Update(rw http.ResponseWriter, r *http.Request) {
	log.Println("Users Update Request")

	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(rw, "Unable to read request data", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	user.ID = id
	err = uh.users.Update(&user, id)
	if err != nil {
		http.Error(rw, "Model not found", http.StatusNotFound)
		return
	}

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err = e.Encode(user)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

// curl -X DELETE localhost:3000/api/users/1
func (uh *UserHandlers) Delete(rw http.ResponseWriter, r *http.Request) {
	log.Println("Users Delete Request")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	var user models.User
	uh.users.Delete(&user, id)

	rw.WriteHeader(http.StatusNoContent)
}
