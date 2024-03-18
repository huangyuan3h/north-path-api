package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const STAGE_KEY = "SST_STAGE"
const PROJECT_STR = "-north-path-api-"

type Client struct {
	TableName *string
	Client    *dynamodb.DynamoDB
}

func New(tableName string) Client {
	client := initDynamo()
	t := getTableName(tableName)
	return Client{TableName: &t, Client: client}
}

func initDynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := dynamodb.New(sess)
	return client
}

func getTableName(tableName string) string {
	var stage string
	if os.Getenv(STAGE_KEY) != "" {
		stage = os.Getenv(STAGE_KEY)
	} else {
		stage = "dev"
	}

	return stage + PROJECT_STR + tableName
}
