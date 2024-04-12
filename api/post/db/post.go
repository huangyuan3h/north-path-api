package db

import (
	"crypto/rand"
	"time"

	db "api.north-path.site/utils/dynamodb"
	"github.com/oklog/ulid"
)

const tableName = "posts"

type Post struct {
	PostId      string   `json:"postId"`
	Email       string   `json:"email"`
	Subject     string   `json:"subject"`
	Content     string   `json:"content"`
	Category    []string `json:"category"`
	CreatedDate string   `json:"createdDate"`
	client      *db.Client
}

type PostMethod interface {
	CreateNew(email, subject, content *string, category *[]string) error
}

func New() PostMethod {
	client := db.New(tableName)

	return Post{client: &client}
}

func (p Post) CreateNew(email, subject, content *string, category *[]string) error {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	post := &Post{
		PostId:      id.String(),
		Email:       *email,
		Subject:     *subject,
		Content:     *content,
		Category:    *category,
		CreatedDate: time.Now().Format(time.RFC3339),
	}

	return p.client.Create(post)
}
