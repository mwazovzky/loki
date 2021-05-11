package main

import (
	"fmt"
	"log"
	"net/http"
)

func mainHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "Hello World!")
}

func main() {
	port := ":3000"

	http.HandleFunc("/", mainHandler)

	log.Println("Listening on port", port)
	http.ListenAndServe(port, nil)
}
