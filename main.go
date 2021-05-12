package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mwazovzky/loki/users"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func homepageHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("Home Page")
	fmt.Fprintln(rw, "Home Page")
}

func userIndex(rw http.ResponseWriter, r *http.Request) {
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

func userShow(rw http.ResponseWriter, r *http.Request) {
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

func userCreate(rw http.ResponseWriter, r *http.Request) {
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

func main() {
	// Configure env variables
	port := ":3000"

	// Setup db connection
	// db = connectDB()

	// Setup routing
	sm := mux.NewRouter()
	sm.HandleFunc("/", homepageHandler)

	apiRouter := sm.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/users", userIndex).Methods(http.MethodGet)
	apiRouter.HandleFunc("/users", userCreate).Methods(http.MethodPost)
	apiRouter.HandleFunc("/users/{id:[0-9]+}", userShow).Methods(http.MethodGet)

	// Configure http server
	server := &http.Server{
		Addr:         port,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// Start http server
	go func() {
		log.Println("Starting http server at", port)
		err := server.ListenAndServe()
		if err != nil {
			log.Println("Error", err)
		}
	}()

	// Gracefully shutdown the server allows to complete current request
	sigChan := make(chan os.Signal)
	// broadcast operating system signals to the channel
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	// wait for the signal
	sig := <-sigChan
	log.Printf("Recieved terminate signal, graceful shutdown, signal: [%s]", sig)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	server.Shutdown(ctx)
}

func connectDB() *gorm.DB {
	dsn := "loki:password@tcp(mysql:3306)/loki"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("ERROR: db connection error")
	}
	return db
}
