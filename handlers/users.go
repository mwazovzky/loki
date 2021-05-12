package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mwazovzky/loki/users"
	"gorm.io/gorm"
)

type UserHandlers struct {
	data *users.Users
}

func NewUserHandlers(db *gorm.DB) *UserHandlers {
	data := users.New(db)
	return &UserHandlers{data}
}

func (uh *UserHandlers) Index(rw http.ResponseWriter, r *http.Request) {
	log.Println("User Index Page")
	fmt.Fprintln(rw, "User Index Page")

	// u := users.New(db)
	// data := u.Find()

	data := []users.User{
		{1, "Mary", "mary@example.com", "password"},
		{2, "Vasy", "vasy@example.com", "password"},
		{3, "Alex", "alex@example.com", "password"},
	}

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err := e.Encode(data)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (uh *UserHandlers) Show(rw http.ResponseWriter, r *http.Request) {
	log.Println("Users Show Page")
	fmt.Fprintln(rw, "User Show Page")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	// u := users.New(db)
	// data := u.Find()

	data := users.User{id, "Test", "test@example.com", "password"}

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err = e.Encode(data)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (uh *UserHandlers) Create(rw http.ResponseWriter, r *http.Request) {
	u := users.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		panic(err)
	}

	rw.Header().Add("Content-Type", "application/json")

	e := json.NewEncoder(rw)
	err = e.Encode(u)
	if err != nil {
		http.Error(rw, "Unable to marshall json", http.StatusInternalServerError)
	}
}
