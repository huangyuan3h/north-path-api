package auth

import (
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
}

const tableName = "auth"

func getTableName() string {
	var stage string
	if os.Getenv("SST_STAGE") != "" {
		stage = os.Getenv("SST_STAGE")
	} else {
		stage = "dev"
	}

	return stage + "-north-path-api-" + tableName
}

func CreateAccount(svc *dynamodb.DynamoDB, email, password *string) error {

	tName := getTableName()

	p := &tName

	auth := Auth{
		Email:    *email,
		Password: *password,
		Status:   "actived",
	}
	av, err := dynamodbattribute.MarshalMap(auth)
	if err != nil {
		return err
	}
	_, err = svc.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: p,
	})

	if err != nil {
		return err
	}

	return nil
}
