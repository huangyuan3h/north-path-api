package auth

import (
	"testing"
)

func TestCreateAccount(t *testing.T) {

	auth := New()

	email := "email1@example.com"
	password := "password"

	err := auth.CreateAccount(&email, &password)

	if err != nil {
		t.Error(err)
	}
}
