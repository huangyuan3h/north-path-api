package auth

import (
	"errors"

	"api.north-path.site/auth/db"
	errs "api.north-path.site/utils/errors"
)

const tableName = "auth"

type Auth struct {
	Email    string `json:"email"`
	Password []byte `json:"password"`
	Status   string `json:"status"`
	client   *db.Client
}

type AuthMethod interface {
	CreateAccount(email, password *string) error
	VerifyLogin(email, password *string) error
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
		return errors.New(errs.AccountAlreadyExists)
	}

	encryptedPassword, err := encrypt([]byte(*password))

	if err != nil {
		return err
	}

	auth := Auth{
		Email:    *email,
		Password: encryptedPassword,
		Status:   "actived",
	}

	return a.client.Create(auth)
}

func (a Auth) VerifyLogin(email, password *string) error {
	item, err := a.client.FindById("email", *email)
	if err != nil {
		return err
	}

	if item == nil {
		return errors.New(errs.AccountNotExist)
	}

	decryptedPassword, err := decrypt(item["password"].B)

	if err != nil {
		return errors.New(errs.PasswordDecryptedError)
	}

	if string(decryptedPassword) != string(*password) {
		return errors.New(errs.PasswordIncorrect)
	}

	return nil
}
