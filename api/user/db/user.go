package db

import (
	"errors"
	"strings"

	db "api.north-path.site/utils/dynamodb"
	errs "api.north-path.site/utils/errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

const tableName = "user"

type User struct {
	Email    string `json:"email" dynamodbav:"email"`
	Avatar   string `json:"avatar" dynamodbav:"avatar"`
	UserName string `json:"userName" dynamodbav:"userName"`
	Bio      string `json:"bio" dynamodbav:"bio"`
	client   *db.Client
}

type UserMethod interface {
	CreateNew(user *User) error
	FindByEmail(email *string) (*User, error)
}

func New() UserMethod {
	client := db.New(tableName)

	return User{client: &client}
}

func (u User) CreateNew(user *User) error {

	return u.client.CreateOrUpdate(user)
}

func (u User) FindByEmail(email *string) (*User, error) {

	item, err := u.client.FindById("email", *email)

	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(item, &u)
	if err != nil {
		return nil, errors.New(errs.UnmarshalError)
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func GetEmailUsername(email string) string {
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		// "@" not found in the email
		return ""
	}

	username := email[:atIndex]
	return username
}
