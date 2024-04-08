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
	FindByEmail(email *string) (*User, error)
}

func New() UserMethod {
	client := db.New(tableName)

	return User{client: &client}
}

func (u User) CreateNew(email *string) error {

	user := &User{
		Email:    *email,
		Avatar:   *email,
		UserName: *email,
		Bio:      *email,
	}

	return u.client.Create(user)
}

func (u User) FindByEmail(email *string) (*User, error) {

	res, err := u.client.FindById("email", *email)

	if err != nil {
		return nil, err
	}

	u.Avatar = *res["avatar"].S
	u.UserName = *res["userName"].S
	u.Bio = *res["bio"].S
	return &u, nil
}
