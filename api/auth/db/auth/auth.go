package auth

import (
	"api.north-path.site/auth/db"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const tableName = "auth"

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
}

type AuthMethod interface {
	CreateAccount(email, password string) error
}

func (a Auth) CreateAccount(email, password *string) error {

	client := db.New(tableName)

	p := client.TableName

	auth := Auth{
		Email:    *email,
		Password: *password,
		Status:   "actived",
	}
	av, err := dynamodbattribute.MarshalMap(auth)
	if err != nil {
		return err
	}
	_, err = client.Client.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: p,
	})

	if err != nil {
		return err
	}

	return nil
}
