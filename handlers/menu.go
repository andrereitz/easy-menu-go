package handlers

import (
	"easy-menu/models"
	"easy-menu/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Menu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	db, _ := utils.Getdb()
	defer db.Close()

	results, err := db.Query("SELECT id, business_logo, business_name, business_color from users WHERE business_url = ?", slug)

	if err != nil {
		http.Error(w, "Error getting user data", http.StatusInternalServerError)
	}

	defer results.Close()

	type BusinessData struct {
		Id            int    `json:"id"`
		BusinessLogo  string `json:"logo"`
		BusinessName  string `json:"name"`
		BusinessColor string `json:"color"`
	}

	var businessData BusinessData

	results.Next()
	err = results.Scan(&businessData.Id, &businessData.BusinessLogo, &businessData.BusinessName, &businessData.BusinessColor)

	if err != nil {
		http.Error(w, "Error reading user data", http.StatusInternalServerError)
	}

	Items, err := db.Query("SELECT id, category, media_id, title, description, price, user FROM Items WHERE user = ?", businessData.Id)

	if err != nil {
		http.Error(w, "Error getting items data", http.StatusInternalServerError)
	}

	defer Items.Close()

	menuItems := make([]models.ItemData, 0)
	for Items.Next() {
		Item := models.ItemData{}
		err := Items.Scan(&Item.Id, &Item.Category, &Item.MediaId, &Item.Title, &Item.Description, &Item.Price, &Item.User)

		if err != nil {
			http.Error(w, "Error reading categories", http.StatusInternalServerError)
			return
		}

		menuItems = append(menuItems, Item)
	}

	generic := utils.ParseSqlNullable(menuItems)

	type Response struct {
		BusinessData BusinessData  `json:"business"`
		Items        []interface{} `json:"items"`
	}

	data := Response{
		BusinessData: businessData,
		Items:        generic,
	}

	respJson, err := json.Marshal(data)

	if err != nil {
		http.Error(w, "Failed to generate json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}
