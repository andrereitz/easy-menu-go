package handlers

import (
	"database/sql"
	"easy-menu/models"
	"easy-menu/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/skip2/go-qrcode"
)

func UserInfo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	db, _ := utils.Getdb()
	defer db.Close()

	row := db.QueryRow("SELECT id, email, business_name, business_url, business_color, business_logo FROM users WHERE id = ?", user)

	var userData models.UserData
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

func UserEditBusiness(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	newUserData := models.UserData{}

	businessName := r.FormValue("business_name")
	newUserData.BusinessName = &businessName
	businessUrl := r.FormValue("business_url")
	newUserData.BusinessUrl = &businessUrl
	businessColor := r.FormValue("business_color")
	newUserData.BusinessColor = &businessColor

	db, _ := utils.Getdb()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE users SET business_name=?, business_url=?, business_color=? WHERE id = ?")

	if err != nil {
		http.Error(w, "Error preparing db operation", http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(newUserData.BusinessName, newUserData.BusinessUrl, newUserData.BusinessColor, user)
	if err != nil {
		http.Error(w, "Error executing db operation", http.StatusInternalServerError)
		return
	}

	if len(businessUrl) > 0 {
		qrpath := "static/qrcodes/" + strconv.Itoa(user) + ".png"
		err := qrcode.WriteFile(businessUrl, qrcode.Medium, 256, qrpath)

		if err != nil {
			fmt.Println("Error generating qrcode")
		}
	}

	response := models.GenericReponse{
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

func UserEditAccount(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	db, _ := utils.Getdb()
	defer db.Close()

	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(password) > 3 {
		passwordhash, err := utils.GetPasswordHash(password)

		if err != nil {
			http.Error(w, "Error generating password hash", http.StatusInternalServerError)
			return
		}

		stmt, err := db.Prepare("UPDATE user SET hash = ? WHERE user = ?")

		if err != nil {
			http.Error(w, "Error preparing db operation", http.StatusInternalServerError)
			return
		}

		defer stmt.Close()

		_, err = stmt.Exec(
			passwordhash,
			user,
		)

		if err != nil {
			http.Error(w, "Error executing db operation", http.StatusInternalServerError)
			return
		}

		// fmt.Println(passwordhash, user)

	}

	stmt, err := db.Prepare("UPDATE user SET email = ? WHERE user = ?")

	if err != nil {
		http.Error(w, "Error preparing db operation", http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		email,
		user,
	)

	if err != nil {
		http.Error(w, "Error executing db operation", http.StatusInternalServerError)
		return
	}

	response := models.GenericReponse{
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

func AddLogo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("logo")

	if err != nil {
		http.Error(w, "Error getting file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	tempFile, err := os.CreateTemp("static/media", "business-logo-*.png")
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

	db, err := utils.Getdb()

	if err != nil {
		http.Error(w, "Error opening db connections", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	stmt, err := db.Prepare("UPDATE users SET business_logo = ? WHERE id = ?")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		tempFile.Name(),
		user,
	)

	if err != nil {
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	response := models.GenericReponse{
		Message: "Image uploaded successfully",
		Status:  "Success",
	}

	respJson, _ := json.Marshal(response)

	w.Header().Set("content-type", "application/json")
	w.Write(respJson)
}

func DeleteLogo(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(int)

	if !ok {
		http.Error(w, "Invalid user!", http.StatusForbidden)
		return
	}

	var logoPath string
	db, _ := utils.Getdb()

	defer db.Close()

	row := db.QueryRow("SELECT business_logo FROM users WHERE id = ?", user)
	err := row.Scan(&logoPath)

	if err != nil {
		fmt.Println(err)
	}

	os.Remove(logoPath)

	stmt, _ := db.Prepare("UPDATE users SET business_logo = ? WHERE id = ?")

	defer stmt.Close()

	stmt.Exec(
		"",
		user,
	)

	response := models.GenericReponse{
		Message: "Logo deleted",
		Status:  "Success",
	}

	respJson, _ := json.Marshal(response)

	w.Header().Set("content-type", "application/json")
	w.Write(respJson)
}
