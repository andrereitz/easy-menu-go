package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type UserData struct {
	Id            *string `json:"id"`
	Email         *string `json:"email"`
	BusinessName  *string `json:"business_name"`
	BusinessUrl   *string `json:"business_url"`
	BusinessColor *string `json:"business_color"`
	BusinessLogo  *string `json:"business_logo"`
}

func UserInfo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(string)

	if r.Method == "GET" {
		if !ok {
			http.Error(w, "User not logged in", http.StatusForbidden)
			return
		}

		path, _ := os.Getwd()
		db, _ := sql.Open("sqlite3", path+"/data/default.db")
		rows, err := db.Query("SELECT id, email, business_name, business_url, business_color, business_logo FROM users WHERE id = ?", user)

		if err != nil {
			http.Error(w, "Failed to get user data", http.StatusInternalServerError)
			return
		}

		var userData UserData

		for rows.Next() {
			err := rows.Scan(&userData.Id, &userData.Email, &userData.BusinessName, &userData.BusinessUrl, &userData.BusinessColor, &userData.BusinessLogo)

			if err != nil {
				fmt.Println(err)
				http.Error(w, "Database Error", http.StatusInternalServerError)
				return
			}
		}

		jsonData, err := json.Marshal(userData)

		if err != nil {
			http.Error(w, "Error while parsing json", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}

}
