package handlers

import (
	"context"
	"easy-menu/models"
	"easy-menu/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("jwt-token")

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		tokenData, verifyErr := utils.VerifyToken(token)

		if verifyErr != nil {
			next.ServeHTTP(w, r)
			return
		}

		id, _ := strconv.Atoi(tokenData["id"])
		ctx := context.WithValue(r.Context(), "user", id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var requestUser models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	user := requestUser.Email
	pass := requestUser.Password

	if user == "" || pass == "" {
		http.Error(w, "please provide all required fields", http.StatusBadRequest)
		return
	}

	db, _ := utils.Getdb()
	rows, err := db.Query("SELECT id, email, hash FROM users WHERE email = ?", user)

	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	var id int
	var email, hash string
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

	compare := utils.ComparePasswordHash([]byte(hash), []byte(pass))

	if compare != nil {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}

	token, err := utils.CreateToken(id)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error generating and encoding token", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "jwt-token",
		Value:    token,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	http.SetCookie(w, cookie)

	response := models.GenericReponse{
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
	_, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "You are not logged in", http.StatusBadRequest)
		return
	}

	c := &http.Cookie{
		Name:     "jwt-token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(w, c)

	response := models.GenericReponse{
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

func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	newUser := models.NewUserData{}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Email and Password is required", http.StatusBadRequest)
		return
	}

	passwordhash, err := utils.GetPasswordHash(password)

	if err != nil {
		http.Error(w, "Error generating password hash", http.StatusBadRequest)
		return
	}

	newUser.Email = email
	newUser.Hash = passwordhash

	db, _ := utils.Getdb()
	stmt, err := db.Prepare("INSERT INTO users (email, hash) VALUES (?, ?)")

	if err != nil {
		http.Error(w, "Error performing db prepare", http.StatusBadRequest)
		return
	}

	_, err = stmt.Exec(newUser.Email, newUser.Hash)

	if err != nil {
		http.Error(w, "Error performing db insertion", http.StatusBadRequest)
		return
	}

	stmt.Close()

	response := models.GenericReponse{
		Message: "User created successfully",
		Status:  "Success",
	}

	respJson, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Error on json marshal", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
