package db

import (
	"testing"

	"api.north-path.site/auth/db/auth"
)

func TestGetClientSanity(t *testing.T) {

	svc := GetClient()

	email := "asd1234@qq.com"
	password := "password"

	err := auth.CreateAccount(svc, &email, &password)

	if err != nil {
		t.Error(err)
	}
}
