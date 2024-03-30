package handlers

import (
	"context"
	"easy-menu/models"
	"easy-menu/utils"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func Items(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	db, _ := utils.Getdb()
	rows, _ := db.Query("SELECT id, category, media_id, title, description, price, user FROM Items WHERE user = ?", user)

	defer db.Close()

	Items := make([]models.ItemData, 0)
	for rows.Next() {
		var Item models.ItemData
		err := rows.Scan(&Item.Id, &Item.Category, &Item.MediaId, &Item.Title, &Item.Description, &Item.Price, &Item.User)
		if err != nil {
			http.Error(w, "Error during row scan", http.StatusInternalServerError)
			return
		}

		Items = append(Items, Item)
	}

	generic := utils.ParseSqlNullable(Items)

	respJson, err := json.Marshal(generic)

	if err != nil {
		http.Error(w, "Error during marshal", http.StatusInternalServerError)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(respJson)
}

func NewItem(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	r.ParseForm()

	var Item models.ItemData
	Item.User = user
	Item.Title = r.FormValue("title")
	Item.Description = utils.NullString(r.FormValue("description"))
	Item.Category = utils.NullIfZero(r.FormValue("category"))
	Item.MediaId = utils.NullIfZero(r.FormValue("media_id"))
	floatVal, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err == nil {
		Item.Price = floatVal
	} else {
		Item.Price = math.NaN()
	}

	db, _ := utils.Getdb()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO items (category, user, media_id, title, description, price) VALUES (?, ?, ?, ?, ?, ?)")

	if err != nil {
		http.Error(w, "Failed to prepare db", http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		Item.Category,
		Item.User,
		Item.MediaId,
		Item.Title,
		Item.Description,
		Item.Price,
	)

	if err != nil {
		http.Error(w, "Failed during dtabase exec", http.StatusInternalServerError)
		return
	}

	response := models.GenericReponse{
		Message: "Item create successfully",
		Status:  "Success",
	}

	respJson, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Error parsing response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func EditItem(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	r.ParseForm()

	var Item models.ItemData
	Item.User = user
	Item.Title = r.FormValue("title")
	Item.Description = utils.NullString(r.FormValue("description"))
	Item.Category = utils.NullIfZero(r.FormValue("category"))
	Item.MediaId = utils.NullIfZero(r.FormValue("media_id"))
	floatVal, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err == nil {
		Item.Price = floatVal
	} else {
		Item.Price = math.NaN()
	}

	db, _ := utils.Getdb()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE items SET category = ?, media_id = ?, title = ?, description = ?, price = ? WHERE id = ? AND user = ?")

	if err != nil {
		http.Error(w, "Error during db prepare", http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		Item.Category,
		Item.MediaId,
		Item.Title,
		Item.Description,
		Item.Price,
		id,
		user,
	)

	if err != nil {
		http.Error(w, "Error during db exec", http.StatusInternalServerError)
		return
	}

	response := models.GenericReponse{
		Message: "Item editted successfully",
		Status:  "Success",
	}

	respJson, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Failed to marshal json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(respJson)
}

func AddItemImage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Item id not provided", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("item_image")

	if err != nil {
		http.Error(w, "Error getting file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	tempFile, err := os.CreateTemp("static/media", "item-image-*.png")
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}

	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading image", http.StatusInternalServerError)
		return
	}

	_, err = tempFile.Write(fileBytes)

	if err != nil {
		http.Error(w, "Error saving image", http.StatusInternalServerError)
		return
	}

	db, _ := utils.Getdb()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		http.Error(w, "Error initiating db transaction", http.StatusInternalServerError)
		return
	}

	ImageData := models.LogoMeta{
		Url:  tempFile.Name(),
		User: user,
	}

	result, err := tx.ExecContext(ctx, "INSERT INTO medias (url, user) VALUES (?, ?)", ImageData.Url, ImageData.User)

	if err != nil {
		http.Error(w, "Transaction failed", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	imageId, errId := result.LastInsertId()

	if errId != nil {
		http.Error(w, "Failed to get image id from result", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	_, errItem := tx.ExecContext(ctx, "UPDATE items SET media_id = ? WHERE id = ? AND user = ?", imageId, id, user)

	if errItem != nil {
		http.Error(w, "Failed to insert image id into item", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	errTx := tx.Commit()

	if errTx != nil {
		http.Error(w, "Error commiting tx", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	response := models.GenericReponse{
		Message: "Item image added",
		Status:  "Success",
	}

	respJson, _ := json.Marshal(response)

	w.Header().Set("content-type", "application/json")
	w.Write(respJson)
}

func RemoveItemImage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Item id not provided", http.StatusBadRequest)
		return
	}

	db, _ := utils.Getdb()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		http.Error(w, "Error initiating db transaction", http.StatusInternalServerError)
		return
	}

	var media_id int
	err = tx.QueryRowContext(ctx, "SELECT media_id FROM items WHERE id = ? AND user = ?", id, user).Scan(&media_id)

	if err != nil {
		http.Error(w, "Invalid image id", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	_, err = tx.ExecContext(ctx, "UPDATE items SET media_id = ? WHERE id = ?", nil, id)

	if err != nil {
		http.Error(w, "Failed to clean item media_id", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	var url string
	err = tx.QueryRowContext(ctx, "SELECT url FROM medias WHERE id = ?", media_id).Scan(&url)

	if err != nil {
		http.Error(w, "Failed to url from media", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	err = os.Remove(url)

	if err != nil {
		http.Error(w, "Failed to remove file from server", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM medias WHERE id = ?", media_id)

	if err != nil {
		http.Error(w, "Failed to delete media form database", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	err = tx.Commit()

	if err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	response := models.GenericReponse{
		Message: "Item image deleted",
		Status:  "Success",
	}

	respJson, _ := json.Marshal(response)

	w.Header().Set("content-type", "application/json")
	w.Write(respJson)
}
