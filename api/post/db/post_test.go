package db

import (
	"testing"
)

func TestCreateNew(t *testing.T) {

	user := New()

	email := "email1@example.com"
	subject := "subject1"
	content := "content1"
	category := []string{"category1"}

	err := user.CreateNew(&email, &subject, &content, &category)

	if err != nil {
		t.Error(err)
	}
}
