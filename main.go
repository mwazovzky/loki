package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mwazovzky/loki/users"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func usersHandler(rw http.ResponseWriter, r *http.Request) {
	db := connectDB()
	u := users.New(db)
	data := u.Find()
	fmt.Println(data)
	fmt.Fprintln(rw, "Hello World!")
}

func main() {
	port := ":3000"

	http.HandleFunc("/users", usersHandler)

	log.Println("Listening on port", port)
	http.ListenAndServe(port, nil)
}

func connectDB() *gorm.DB {
	dsn := "loki:password@tcp(mysql:3306)/loki"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("ERROR: db connection error")
	}
	return db
}
