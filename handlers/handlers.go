package handlers

import (
	"easy-menu/models"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["index"] = template.Must(template.ParseFiles("templates/layouts/clean.html", "templates/shared/head.html", "templates/pages/index.html"))
	templates["dashboard"] = template.Must(template.ParseFiles("templates/layouts/default.html", "templates/shared/head.html", "templates/shared/navbar.html", "templates/pages/index.html"))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Println("Index handler")

	tmpl := templates["index"]
	tmpl.Execute(w, models.Example{Title: "Some title", Description: "Some description", Session: 1})
}

func ItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Item: %v\n", vars["id"])
}
