package utils

import (
	"database/sql"
	"easy-menu/models"
	"errors"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte
var hashKey []byte
var blockKey []byte
var s *securecookie.SecureCookie

func init() {
	godotenv.Load()

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	hashKey = []byte(os.Getenv("HASH_KEY"))
	blockKey = []byte(os.Getenv("BLOCK_KEY"))
	s = securecookie.New(hashKey, blockKey)
}

func CreateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		"id":         strconv.Itoa(userID),
		"exp":        time.Now().Add(time.Minute * 1440).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	encoded, cookieErr := s.Encode("jwt-token", tokenString)

	if cookieErr != nil {
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
	err := bcrypt.CompareHashAndPassword(storedPassword, password)
	if err != nil {
		return errors.New("password does not match")
	}

	return nil
}

func NullIfZero(s string) sql.NullInt64 {
	n, err := strconv.Atoi(s)
	if err != nil || n == 0 {
		return sql.NullInt64{}
	}

	return sql.NullInt64{
		Int64: int64(n),
		Valid: true,
	}
}

func NullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func ParseSqlNullable(Items []models.ItemData) []interface{} {
	GenericItems := make([]interface{}, 0)

	for _, item := range Items {
		GenericItem := make(map[string]interface{})
		t := reflect.TypeOf(item)

		for i := 0; i < t.NumField(); i++ {
			field := reflect.ValueOf(item).Field(i)
			fieldType := field.Type()
			fieldName := t.Field(i)

			if fieldType.String() == "sql.NullInt64" {
				nullInt64 := field.Interface().(sql.NullInt64)
				if nullInt64.Valid {
					GenericItem[fieldName.Name] = nullInt64.Int64
				}
			} else if fieldType.String() == "sql.NullString" {
				nullString := field.Interface().(sql.NullString)
				if nullString.Valid {
					GenericItem[fieldName.Name] = nullString.String
				}
			} else {
				value := field.Interface()

				GenericItem[fieldName.Name] = value
			}
		}

		GenericItems = append(GenericItems, GenericItem)
	}

	return GenericItems
}
