package handlers

import (
	"database/sql"
	"easy-menu/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func UserInfo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	if r.Method == "GET" {

		db, _ := utils.Getdb()
		defer db.Close()

		row := db.QueryRow("SELECT id, email, business_name, business_url, business_color, business_logo FROM users WHERE id = ?", user)

		var userData UserData
		err := row.Scan(&userData.Id, &userData.Email, &userData.BusinessName, &userData.BusinessUrl, &userData.BusinessColor, &userData.BusinessLogo)

		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Error fetching category", http.StatusNotFound)
			return
		}

		jsonData, err := json.Marshal(userData)

		if err != nil {
			http.Error(w, "Error while parsing json", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}

	if r.Method == "POST" {
		r.ParseForm()
		newUserData := UserData{}
		newUserData.Email = r.FormValue("email")

		businessName := r.FormValue("business_name")
		newUserData.BusinessName = &businessName
		businessUrl := r.FormValue("business_url")
		newUserData.BusinessUrl = &businessUrl
		businessColor := r.FormValue("business_color")
		newUserData.BusinessColor = &businessColor
		businessLogo := r.FormValue("business_logo")
		newUserData.BusinessLogo = &businessLogo

		db, _ := utils.Getdb()

		stmt, err := db.Prepare("UPDATE users SET email=?, business_name=?, business_url=?, business_color=?, business_logo=? WHERE id = ?")
		if err != nil {
			http.Error(w, "Error preparing db operation", http.StatusInternalServerError)
			return
		}

		defer stmt.Close()

		_, err = stmt.Exec(newUserData.Email, newUserData.BusinessName, newUserData.BusinessUrl, newUserData.BusinessColor, newUserData.BusinessLogo, user)
		if err != nil {
			fmt.Println("Error executing db operation", err)

			http.Error(w, "Error executing db operation", http.StatusInternalServerError)
			return
		}

		response := GenericReponse{
			Message: "User updated successfully",
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
}
