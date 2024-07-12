package handlers

import (
	"context"
	"database/sql"
	"easy-menu/models"
	"easy-menu/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Category(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	db, _ := utils.Getdb()
	defer db.Close()

	row := db.QueryRow("SELECT id, user, title FROM categories WHERE id = ? AND user = ?", id, user)

	var category models.CategoryData
	err := row.Scan(&category.Id, &category.User, &category.Title)

	if err == sql.ErrNoRows {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching category", http.StatusNotFound)
		return
	}

	response := models.DataReponse{
		Data: category,
	}

	respJson, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Error during json marshal", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func Categories(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	db, _ := utils.Getdb()
	defer db.Close()

	rows, err := db.Query("SELECT id, user, title FROM categories WHERE user = ?", user)

	if err != nil {
		http.Error(w, "error seelct", http.StatusForbidden)
		return
	}

	defer rows.Close()

	categories := make([]models.CategoryData, 0)
	for rows.Next() {
		category := models.CategoryData{}
		err := rows.Scan(&category.Id, &category.User, &category.Title)

		if err != nil {
			http.Error(w, "Error reading categories", http.StatusInternalServerError)
			return
		}

		categories = append(categories, category)
	}

	respJson, err := json.Marshal(categories)

	if err != nil {
		http.Error(w, "Failed to generate json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func NewCategory(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	title := r.FormValue("title")

	db, _ := utils.Getdb()
	stmt, _ := db.Prepare("INSERT INTO categories (user, title) VALUES (?,?)")
	_, err := stmt.Exec(user, title)

	if err != nil {
		http.Error(w, "Error inserting new category", http.StatusInternalServerError)
		return
	}

	response := models.GenericReponse{
		Message: "Category added!",
		Status:  "Success",
	}

	respJson, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func EditCategory(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	db, _ := utils.Getdb()
	row := db.QueryRow("SELECT * FROM categories WHERE id = ? AND user = ?", id, user)

	var category models.CategoryData
	err := row.Scan(&category.Id, &category.User, &category.Title)

	if err != nil {
		http.Error(w, "Database select error", http.StatusForbidden)
		return
	}

	if category.User != user || err != nil {
		http.Error(w, "Can't edit this category", http.StatusForbidden)
		return
	}

	newTitle := r.FormValue("title")
	stmt, err := db.Prepare("UPDATE categories SET title = ? WHERE id = ?")

	if err != nil {
		http.Error(w, "Database prepare error", http.StatusForbidden)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(newTitle, id)

	if err != nil {
		http.Error(w, "Database exec error", http.StatusForbidden)
		return
	}

	response := models.GenericReponse{
		Message: "Editted successfully",
		Status:  "Success",
	}

	respJson, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	db, _ := utils.Getdb()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		http.Error(w, "Error initiating db transaction", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	_, categoryErr := tx.ExecContext(ctx, "DELETE FROM categories WHERE id = ? AND user = ?", id, user)

	if categoryErr != nil {
		http.Error(w, "Category transaction failed", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	_, errItems := tx.ExecContext(ctx, "UPDATE items SET category = ? WHERE category = ? AND user = ?", nil, id, user)

	if errItems != nil {
		http.Error(w, "Category transaction failed to execute on items", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	errTx := tx.Commit()

	if errTx != nil {
		http.Error(w, "Error commiting trasaction", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	response := models.GenericReponse{
		Message: "Deleted successfully",
		Status:  "Success",
	}

	respJson, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}
