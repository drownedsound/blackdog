package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /application/new", newApplication)
	mux.HandleFunc("POST /application/new", postNewApplication)
	mux.HandleFunc("GET /application/{id}/view", viewApplication)
	mux.HandleFunc("GET /application/{id}", getApplicationData)
	// TODO: Create handler for different application versions
	// For example, /application/100/version/1 or
	// /application/100/version/latest
	mux.HandleFunc("GET /application/{id}/edit", editApplication)

	log.Println("Startng HTTP server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
