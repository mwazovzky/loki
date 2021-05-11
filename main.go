package main

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func mainHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Hello World!")

	dsn := "loki:password@tcp(mysql:3306)/loki"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("ERROR: db connection error", err)
	}

	var users []User
	result := db.Find(&users)

	fmt.Println(users)
	fmt.Println("RowsAffected:", result.RowsAffected, "ERROR", result.Error)
}

func main() {
	port := ":3000"

	http.HandleFunc("/", mainHandler)

	log.Println("Listening on port", port)
	http.ListenAndServe(port, nil)
}
