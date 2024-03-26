package handlers

import (
	"context"
	"database/sql"
	"easy-menu/models"
	"easy-menu/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["index"] = template.Must(template.ParseFiles("templates/layouts/clean.html", "templates/shared/head.html", "templates/pages/index.html"))
	templates["dashboard"] = template.Must(template.ParseFiles("templates/layouts/default.html", "templates/shared/head.html", "templates/shared/navbar.html", "templates/pages/dashboard.html"))
	templates["login"] = template.Must(template.ParseFiles("templates/layouts/clean.html", "templates/shared/head.html", "templates/pages/login.html"))
}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("jwt-token")

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		tokenData, verifyErr := utils.VerifyToken(token)

		if verifyErr != nil {
			fmt.Println("Verify error: ", verifyErr)

			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tokenData["id"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	var pageData = map[string]string{}
	user, ok := r.Context().Value("user").(string)

	if ok {
		pageData["user"] = user
	}

	tmpl := templates["index"]
	tmpl.Execute(w, models.PageData{Data: pageData})
}

type GenericReponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := r.FormValue("email")
	pass := r.FormValue("password")

	path, _ := os.Getwd()
	db, _ := sql.Open("sqlite3", path+"/data/default.db")
	rows, err := db.Query("SELECT id, email, hash FROM users WHERE email = ?", user)

	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	var id, email, hash string
	for rows.Next() {
		err := rows.Scan(&id, &email, &hash)

		if err != nil {
			http.Error(w, "Database Error", http.StatusInternalServerError)
		}
	}

	db.Close()

	if email == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if pass != hash {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	intId, _ := strconv.Atoi(id)
	token, err := utils.CreateToken(intId)

	if err != nil {
		http.Error(w, "Error generating and encoding token", http.StatusInternalServerError)
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

	response := GenericReponse{
		Message: "Logged in successfully",
		Status:  "Success",
	}

	respJson, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user").(string)

	if !ok {
		http.Error(w, "You are not logged in", http.StatusBadRequest)
		return
	}

	c := &http.Cookie{
		Name:    "jwt-token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}

	http.SetCookie(w, c)

	response := GenericReponse{
		Message: "Logged out successfully",
		Status:  "Success",
	}

	respJson, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	var pageData = map[string]string{}
	user, ok := r.Context().Value("user").(string)

	if ok {
		pageData["user"] = user
	}
}

func ItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Item: %v\n", vars["id"])
}
