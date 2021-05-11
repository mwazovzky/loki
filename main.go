package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/mwazovzky/loki/users"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func homepageHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("Home Page")
	fmt.Fprintln(rw, "Home Page")
}

func usersHandler(rw http.ResponseWriter, r *http.Request) {
	log.Println("Users Index Page")
	db := connectDB()
	u := users.New(db)
	data := u.Find()
	fmt.Println(data)
	fmt.Fprintln(rw, "Users Page")
}

func main() {
	// Configure env variables
	port := ":3000"

	// Setup routing
	sm := http.NewServeMux()
	sm.HandleFunc("/", homepageHandler)
	sm.HandleFunc("/users", usersHandler)

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
