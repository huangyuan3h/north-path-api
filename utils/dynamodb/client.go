package dynamodb

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const STAGE_KEY = "SST_STAGE"
const PROJECT_STR = "-north-path-api-"

type Client struct {
	TableName *string
	Client    *dynamodb.DynamoDB
}

type ClientBaseMethod interface {
	CreateOrUpdate(in interface{}) error
	FindById(keyName, id string) (map[string]*dynamodb.AttributeValue, error)
	DeleteById(keyName, id string) error
}

func (c Client) CreateOrUpdate(in interface{}) error {
	av, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return err
	}
	_, err = c.Client.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: c.TableName,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c Client) FindById(keyName, id string) (map[string]*dynamodb.AttributeValue, error) {
	if keyName == "" {
		return nil, errors.New("keyname is required")
	}
	key := map[string]*dynamodb.AttributeValue{
		keyName: {
			S: aws.String(id),
		},
	}

	return GetItem(c.Client, *c.TableName, key)
}

func (c Client) DeleteById(keyName, id string) error {
	if keyName == "" {
		return errors.New("keyname is required")
	}
	key := map[string]*dynamodb.AttributeValue{
		keyName: {
			S: aws.String(id),
		},
	}

	return DeleteItem(c.Client, *c.TableName, key)
}

func GetItem(svc *dynamodb.DynamoDB, tableName string, key map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	result, err := svc.GetItem(input)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

func DeleteItem(svc *dynamodb.DynamoDB, tableName string, key map[string]*dynamodb.AttributeValue) error {

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
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
