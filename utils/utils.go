package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("sometestsecretkey765")
var hashKey = []byte("sometesthashkey-securecookie123")
var blockKey = []byte("sometestblockkey-securecookie123")
var s = securecookie.New(hashKey, blockKey)

func CreateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"id":         strconv.Itoa(userID),
		"exp":        time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	encoded, cookieErr := s.Encode("jwt-token", tokenString)

	if cookieErr != nil {
		fmt.Println(cookieErr)
		return "", cookieErr
	}

	return encoded, nil
}

func VerifyToken(token *http.Cookie) (map[string]string, error) {
	var decoded string
	value := token.Value
	err := s.Decode("jwt-token", value, &decoded)

	if err != nil {
		return map[string]string{}, err
	}

	claims := jwt.MapClaims{}
	_, tokenErr := jwt.ParseWithClaims(decoded, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	fmt.Println(claims["id"])
	finalMap := map[string]string{}

	for key, val := range claims {
		if str, ok := val.(string); ok {
			finalMap[key] = str
		} else {
			finalMap[key] = string(str)
		}
	}

	return finalMap, tokenErr
}

func Getdb() (*sql.DB, error) {
	path, err := os.Getwd()

	if err != nil {
		return nil, errors.New("falied to open current working directory")
	}

	db, err := sql.Open("sqlite3", path+"/data/default.db")

	if err != nil {
		return nil, errors.New("falied to open connection to database")
	}

	return db, nil
}

func GetPasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	hashedPasswordString := string(hashedPassword)

	return hashedPasswordString, nil
}

func ComparePasswordHash(storedPassword []byte, password []byte) error {
	// storedPasswordhash, err := GetPasswordHash(string(dbpassword))

	// if err != nil {
	// 	return errors.New("failed to hash provided password")
	// }

	err := bcrypt.CompareHashAndPassword(storedPassword, password)
	if err != nil {
		return errors.New("password does not match")
	}

	return nil
}
