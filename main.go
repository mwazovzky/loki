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
	"github.com/joho/godotenv"
	"github.com/mwazovzky/loki/handlers"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var port string
var db *gorm.DB

func init() {
	godotenv.Load()
	port = fmt.Sprintf(":%s", os.Getenv("PORT"))
	db = connectDB()
}

func main() {
	// Setup routing
	sm := mux.NewRouter()

	userHandlers := handlers.NewUserHandlers(db)
	apiRouter := sm.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/users", userHandlers.Index).Methods(http.MethodGet)
	apiRouter.HandleFunc("/users", userHandlers.Create).Methods(http.MethodPost)
	apiRouter.HandleFunc("/users/{id:[0-9]+}", userHandlers.Show).Methods(http.MethodGet)
	apiRouter.HandleFunc("/users/{id:[0-9]+}", userHandlers.Delete).Methods(http.MethodDelete)
	apiRouter.HandleFunc("/users/{id:[0-9]+}/update", userHandlers.Update).Methods(http.MethodPut)

	pagesHandlers := handlers.NewPagesHandlers()
	pagesRouter := sm.Methods(http.MethodGet).Subrouter()

	pagesRouter.HandleFunc("/", pagesHandlers.Home)

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
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("ERROR: db connection error")
	}

	return db
}
