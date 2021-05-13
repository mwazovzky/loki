package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type PagesHandlers struct {
	//
}

func NewPagesHandlers() *PagesHandlers {
	return &PagesHandlers{}
}

func (*PagesHandlers) Home(rw http.ResponseWriter, r *http.Request) {
	log.Println("Home Page")
	fmt.Fprintln(rw, "Welcome to Home Page")
}
