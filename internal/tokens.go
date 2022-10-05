package internal

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(login string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(72 * time.Hour).Unix()
	claims["user"] = login
	tokenString, err := token.SignedString([]byte(os.Getenv("secret")))
	if err != nil {
		log.Print(err)
		return "Signing Error", err
	}

	return tokenString, nil
}
