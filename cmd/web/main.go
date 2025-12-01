package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	type config struct {
		addr      string
		staticDir string
	}

	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
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

	log.Println("Startng HTTP server on " + cfg.addr)
	err := http.ListenAndServe(cfg.addr, mux)
	log.Fatal(err)
}
