package db

import (
	"strings"

	db "api.north-path.site/utils/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
		UserName: getEmailUsername(*email),
	}

	return u.client.Create(user)
}

func (u User) FindByEmail(email *string) (*User, error) {

	res, err := u.client.FindById("email", *email)

	if err != nil {
		return nil, err
	}

	mappedU, err := map2User(res)

	if err != nil {
		return nil, err
	}

	return mappedU, nil
}

func map2User(res map[string]*dynamodb.AttributeValue) (*User, error) {
	user := &User{}

	if err := assignStringAttribute(res, "email", &user.Email); err != nil {
		return nil, err
	}

	if err := assignStringAttribute(res, "avatar", &user.Avatar); err != nil {
		return nil, err
	}

	if err := assignStringAttribute(res, "userName", &user.UserName); err != nil {
		return nil, err
	}

	if err := assignStringAttribute(res, "bio", &user.Bio); err != nil {
		return nil, err
	}

	return user, nil
}

func assignStringAttribute(res map[string]*dynamodb.AttributeValue, key string, target *string) error {
	if res[key] == nil || res[key].S == nil {
		*target = ""
	} else {
		*target = *res[key].S
	}
	return nil
}

func getEmailUsername(email string) string {
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		// "@" not found in the email
		return ""
	}

	username := email[:atIndex]
	return username
}
