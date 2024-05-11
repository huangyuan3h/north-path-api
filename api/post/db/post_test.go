package db

import (
	"testing"

	types "api.north-path.site/post/types"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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

	const id = "01HXGY4M2327R3TW5TY6H1BT6K"

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

	err := post.DeleteById("01HXEE00422CT5ARXN05G56RSQ")
	if err != nil {
		t.Error(err)
	}
}

func TestSearch(t *testing.T) {
	post := New()

	p := types.SearchKeys{PostId: "01HXJVRDN4WK6T6TF6PZYRC6FY", UpdatedDate: "2024-05-11T11:39:44+08:00"}

	a, err := attributevalue.MarshalMap(p)

	if err != nil {
		t.Error(err)
	}

	item, nextToken, err := post.Search(1, a, "")
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
