package db

import (
	"testing"

	"api.north-path.site/user/types"
)

func TestCreateNew(t *testing.T) {

	user := New()
	email := "email2@example.com"
	u := &types.User{
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

	email := "email2@example.com"

	res, err := user.FindByEmail(&email)

	if err != nil {
		t.Error(err)
	}
	if res.Email != email {
		t.Error("email not equal")
	}
}
