package helper

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(id int64, email, name string) string {
	claims := jwt.MapClaims{
		"id":       id,
		"name":     name,
		"email":    email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"exp_date": time.Now().Add(time.Hour * 24),
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, _ := parseToken.SignedString([]byte(os.Getenv("SECRET_KEY")))

	return signedToken
}

func VerifyToken(tokenString string) (interface{}, error) {
	errResponse := errors.New("Token-Invalid")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errResponse
		}

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, errResponse
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errResponse
	}

	return token.Claims.(jwt.MapClaims), nil
}
