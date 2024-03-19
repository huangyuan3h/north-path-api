package auth

import (
	"errors"

	"api.north-path.site/auth/db"
)

const tableName = "auth"

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
	client   *db.Client
}

type AuthMethod interface {
	CreateAccount(email, password string) error
}

func New() Auth {
	client := db.New(tableName)

	return Auth{client: &client}
}

func (a Auth) CreateAccount(email, password *string) error {

	item, err := a.client.FindById("email", *email)

	if err != nil {
		return err
	}

	if item != nil {
		return errors.New("the account already exists")
	}

	auth := Auth{
		Email:    *email,
		Password: *password,
		Status:   "actived",
	}

	return a.client.Create(auth)
}
