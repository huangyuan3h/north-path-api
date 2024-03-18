package auth

import (
	"testing"
)

func TestCreateAccount(t *testing.T) {

	auth := Auth{}

	email := "email@example.com"
	password := "password"

	err := auth.CreateAccount(&email, &password)

	if err != nil {
		t.Error(err)
	}
}
