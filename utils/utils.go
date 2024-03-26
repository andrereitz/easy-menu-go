package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/securecookie"
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
