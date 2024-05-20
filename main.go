package main

import (
	"easy-menu/handlers"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var dir string

	flag.StringVar(&dir, "dir", "./static", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	r := mux.NewRouter()
	r.Use(handlers.Authorization)
	r.Use(handlers.CorsMiddleware)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/authorize", handlers.Authorize).Methods("GET")
	r.HandleFunc("/logout", handlers.Logout).Methods("POST")
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/user", handlers.UserInfo).Methods("GET", "POST")
	r.HandleFunc("/user/logo/add", handlers.AddLogo).Methods("POST")
	r.HandleFunc("/user/logo/delete", handlers.DeleteLogo).Methods("POST")
	r.HandleFunc("/category/new", handlers.NewCategory).Methods("POST")
	r.HandleFunc("/category/all", handlers.Categories).Methods("GET")
	r.HandleFunc("/category/{id}", handlers.Category).Methods("GET")
	r.HandleFunc("/category/edit/{id}", handlers.EditCategory).Methods("POST")
	r.HandleFunc("/category/delete/{id}", handlers.DeleteCategory).Methods("POST")
	r.HandleFunc("/item/all", handlers.Items).Methods("GET")
	r.HandleFunc("/item/new", handlers.NewItem).Methods("POST")
	r.HandleFunc("/item/edit/{id}", handlers.EditItem).Methods("POST")
	r.HandleFunc("/item/delete/{id}", handlers.DeleteItem).Methods("POST")
	r.HandleFunc("/item/image/{id}", handlers.GetItemImage).Methods("GET")
	r.HandleFunc("/item/image/add/{id}", handlers.AddItemImage).Methods("POST")
	r.HandleFunc("/item/image/remove/{id}", handlers.RemoveItemImage).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
