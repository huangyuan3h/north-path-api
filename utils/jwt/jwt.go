package jwt

import (
	"time"

	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Email     string
	Username  string
	ExpiresAt time.Time
	Issuer    string
}

const SECRET_KEY = "JWT_SECRET"

func getKey() ([]byte, error) {

	if os.Getenv(SECRET_KEY) == "" {
		return nil, errors.New("secret is not set")
	}
	return []byte(os.Getenv(SECRET_KEY)), nil
}

func CreateToken(in map[string]interface{}) (string, error) {
	secret, err := getKey()
	if err != nil {
		return "", err
	}
	claim := make(jwt.MapClaims, len(in)+2) // add issuer and expires
	claim["iss"] = "http://north-path.site"
	claim["exp"] = time.Now().Add(time.Hour * 24).Unix()
	for key, value := range in {
		claim[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) error {

	secret, err := getKey()
	if err != nil {
		return err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("token is invalid")
	}

	return nil
}
