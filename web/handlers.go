package web

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"

	"github.com/drownedsound/blackdog/common"
)

func HomeHandler(svc *common.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "Go")

		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			"./ui/html/partials/sidebar.tmpl",
			"./ui/html/pages/home.tmpl",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			svc.ServerError(w, r, err)
			return
		}

		svc.Logger.Debug(
			"Rendering home page",
			slog.String("method", "GET"),
			slog.String("path", "/"),
		)

		err = ts.ExecuteTemplate(w, "base", nil)
		if err != nil {
			svc.ServerError(w, r, err)
		}
	}
}

func NewApplicationHandler(svc *common.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Please create a new application 🗒️")
	}
}

func PostNewApplicationHandler(svc *common.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "You just created a new application 🗒️")
	}
}

func ViewApplicationHandler(svc *common.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
		io.WriteString(w, "You are now viewing application ")
		io.WriteString(w, strconv.Itoa(id))
	}
}

func EditApplicationHandler(svc *common.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
		io.WriteString(w, "You are now editing application ")
		io.WriteString(w, strconv.Itoa(id))
	}
}

func GetApplicationDataHandler(svc *common.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{ id=%d, name="Drowned Sound"}`, id)
	}
}
