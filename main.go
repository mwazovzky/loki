package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/mwazovzky/loki/handlers"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func homepageHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("Home Page")
	fmt.Fprintln(rw, "Home Page")
}

func main() {
	// Configure env variables
	port := ":3000"

	// Setup db connection
	db = connectDB()

	// Setup routing
	sm := mux.NewRouter()
	sm.HandleFunc("/", homepageHandler)

	userHandlers := handlers.NewUserHandlers(db)
	apiRouter := sm.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/users", userHandlers.Index).Methods(http.MethodGet)
	apiRouter.HandleFunc("/users", userHandlers.Create).Methods(http.MethodPost)
	apiRouter.HandleFunc("/users/{id:[0-9]+}", userHandlers.Show).Methods(http.MethodGet)

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
