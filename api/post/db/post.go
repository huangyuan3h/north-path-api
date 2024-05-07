package db

import (
	"crypto/rand"
	"time"

	db "api.north-path.site/utils/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/oklog/ulid"
)

const tableName = "posts"

type Post struct {
	PostId      string   `json:"postId"`
	Email       string   `json:"email"`
	Subject     string   `json:"subject"`
	Content     string   `json:"content"`
	Categories  []string `json:"categories"`
	Images      []string `json:"images"`
	CreatedDate string   `json:"createdDate"`
	UpdatedDate string   `json:"updatedDate"`
	client      *db.Client
}

type PostMethod interface {
	CreateNew(email, subject, content *string, images, categories *[]string) (Post, error)
	FindById(id string) (*Post, error)
	DeleteById(id string) error
}

func New() PostMethod {
	client := db.New(tableName)

	return Post{client: &client}
}

func (p Post) CreateNew(email, subject, content *string, images, categories *[]string) (Post, error) {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	post := &Post{
		PostId:      id.String(),
		Email:       *email,
		Subject:     *subject,
		Content:     *content,
		Categories:  *categories,
		Images:      *images,
		CreatedDate: time.Now().Format(time.RFC3339),
		UpdatedDate: time.Now().Format(time.RFC3339),
	}

	return *post, p.client.CreateOrUpdate(post)
}

func map2Post(postMap map[string]*dynamodb.AttributeValue) (*Post, error) {
	post := &Post{}
	err := dynamodbattribute.UnmarshalMap(postMap, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p Post) FindById(id string) (*Post, error) {
	postMap, err := p.client.FindById("postId", id)
	if err != nil {
		return nil, err
	}
	return map2Post(postMap)
}

func (p Post) DeleteById(id string) error {
	return p.client.DeleteById("postId", id)
}
