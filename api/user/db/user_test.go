package db

import (
	"testing"
)

func TestCreateNew(t *testing.T) {

	user := New()
	email := "email1@example.com"
	u := &User{
		Email:    email,
		UserName: GetEmailUsername(email),
	}

	err := user.CreateNew(u)

	if err != nil {
		t.Error(err)
	}
}

func TestFindByEmail(t *testing.T) {
	user := New()

	email := "email1@example.com"

	res, err := user.FindByEmail(&email)

	if err != nil {
		t.Error(err)
	}
	if res.Email != email {
		t.Error("email not equal")
	}
}
