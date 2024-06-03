package db

import (
	"errors"
	"strings"

	"api.north-path.site/user/types"
	db "api.north-path.site/utils/dynamodb"
	errs "api.north-path.site/utils/errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

const tableName = "user"

type User struct {
	types.User
	client *db.Client
}

type UserMethod interface {
	CreateNew(user *types.User) error
	FindByEmail(email *string) (*types.User, error)
}

func New() UserMethod {
	client := db.New(tableName)

	return User{client: &client}
}

func (u User) CreateNew(user *types.User) error {

	return u.client.CreateOrUpdate(user)
}

func (u User) FindByEmail(email *string) (*types.User, error) {

	item, err := u.client.FindById("email", *email)

	if err != nil {
		return nil, err
	}

	user := types.User{}

	err = attributevalue.UnmarshalMap(item, &user)
	if err != nil {
		return nil, errors.New(errs.UnmarshalError)
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
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
