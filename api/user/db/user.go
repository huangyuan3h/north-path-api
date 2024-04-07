package db

import (
	db "api.north-path.site/utils/dynamodb"
)

const tableName = "user"

type User struct {
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	UserName string `json:"userName"`
	Bio      string `json:"bio"`
	client   *db.Client
}

type UserMethod interface {
	CreateNew(email *string) error
}

func New() UserMethod {
	client := db.New(tableName)

	return User{client: &client}
}

func (a User) CreateNew(email *string) error {

	user := &User{
		Email: *email,
	}

	return a.client.Create(user)
}
