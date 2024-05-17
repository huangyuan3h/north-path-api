package db

import (
	"testing"
)

func TestCreateNew(t *testing.T) {

	post := New()

	email := "email1@example.com"
	subject := "subject1"
	content := "content1"
	category := "category1"
	location := "location1"
	topics := []string{"topic1"}
	images := []string{"image1"}

	_, err := post.CreateNew(&email, &subject, &content, &category, &location, &images, &topics)

	if err != nil {
		t.Error(err)
	}
}

func TestFindById(t *testing.T) {
	post := New()

	const id = "01HY2S04E8H8JD6DSQ2E6RWBRW"

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

	err := post.DeleteById("01HY2S04E8H8JD6DSQ2E6RWBRW")
	if err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	post := New()

	item, nextToken, err := post.Search(2, "", "category1")
	if err != nil {
		t.Error(err)
	}
	if len(item) == 0 {
		t.Error("shoule not empty")
	}

	if nextToken == nil {
		t.Error("nextToken should not be nil")
	}

}
