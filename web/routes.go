package web

import (
	"net/http"

	"github.com/drownedsound/blackdog/common"
)

func Routes(dir string, svc *common.Services) *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(dir))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", HomeHandler(svc))
	mux.HandleFunc("GET /application/new", NewApplicationHandler(svc))
	mux.HandleFunc("POST /application/new", PostNewApplicationHandler(svc))
	mux.HandleFunc("GET /application/{id}/view", ViewApplicationHandler(svc))
	mux.HandleFunc("GET /application/{id}", GetApplicationDataHandler(svc))
	// TODO: Create handler for different application versions
	// For example, /application/100/version/1 or
	// /application/100/version/latest
	mux.HandleFunc("GET /application/{id}/edit", EditApplicationHandler(svc))

	return mux
}
