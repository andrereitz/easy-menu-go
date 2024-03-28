package handlers

import (
	"easy-menu/models"
	"easy-menu/utils"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
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

	response := GenericReponse{
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
