package db

import (
	"testing"
)

func TestCreateNew(t *testing.T) {

	post := New()

	email := "email1@example.com"
	subject := "subject1"
	content := "content1"
	category := []string{"category1"}
	images := []string{"image1"}

	err := post.CreateNew(&email, &subject, &content, &images, &category)

	if err != nil {
		t.Error(err)
	}
}

func TestDeleteById(t *testing.T) {
	post := New()

	err := post.DeleteById("01HV9FKH5H8NQYYPTKCKVAZD3E")
	if err != nil {
		t.Error(err)
	}
}
