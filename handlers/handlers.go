package handlers

import (
	"database/sql"
	"easy-menu/models"
	"easy-menu/utils"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["index"] = template.Must(template.ParseFiles("templates/layouts/clean.html", "templates/shared/head.html", "templates/pages/index.html"))
	templates["dashboard"] = template.Must(template.ParseFiles("templates/layouts/default.html", "templates/shared/head.html", "templates/shared/navbar.html", "templates/pages/index.html"))
	templates["login"] = template.Must(template.ParseFiles("templates/layouts/clean.html", "templates/shared/head.html", "templates/pages/login.html"))
}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		token, err := r.Cookie("jwt-token")

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		valid, verifyErr := utils.VerifyToken(token)

		if verifyErr != nil {
			fmt.Println(verifyErr)

			next.ServeHTTP(w, r)
			return
		}

		if valid {
			context.Set(r, "user", token)
		}

		next.ServeHTTP(w, r)
	})
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Println("Index handler")

	tmpl := templates["index"]
	tmpl.Execute(w, models.Example{Title: "Some title", Description: "Some description", Session: 1})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		user := r.FormValue("email")
		pass := r.FormValue("password")

		fmt.Println("user and pass", user, pass)

		path, _ := os.Getwd()
		db, _ := sql.Open("sqlite3", path+"/data/default.db")
		rows, err := db.Query("SELECT id, email, hash FROM users WHERE email = ?", user)

		if err != nil {
			fmt.Println("Error executing query")
		}

		var id, email, hash string
		for rows.Next() {
			err := rows.Scan(&id, &email, &hash)

			if err != nil {
				fmt.Println(err)
			}
		}

		if email == "" {
			fmt.Println("User not found")

			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		fmt.Println(email, hash, "---> email and hash")

		if pass != hash {
			fmt.Println("Invalid password")

			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		intId, _ := strconv.Atoi(id)

		token, err := utils.CreateToken(intId)

		if err != nil {
			fmt.Println("Error generating and encoding token")
			return
		}

		cookie := &http.Cookie{
			Name:     "jwt-token",
			Value:    token,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}

	if r.Method == "GET" {
		tmpl := templates["login"]
		tmpl.Execute(w, models.Example{Title: "Some title", Description: "Some description", Session: 1})
	}
}

func ItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Item: %v\n", vars["id"])
}
