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

	_, err := post.CreateNew(&email, &subject, &content, &images, &category)

	if err != nil {
		t.Error(err)
	}
}

func TestFindById(t *testing.T) {
	post := New()

	const id = "01HXDZY7K8X58EMP171GT32KDQ"

	item, err := post.FindById(id)
	if err != nil {
		t.Error(err)
	}

	if item.PostId != id {
		t.Error("post id not equal")
	}
}

func TestDeleteById(t *testing.T) {
	post := New()

	err := post.DeleteById("01HV9FKH5H8NQYYPTKCKVAZD3E")
	if err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	post := New()

	item, err := post.Search()
	if err != nil {
		t.Error(err)
	}
	t.Error(item)
}
