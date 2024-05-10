package db

import (
	"crypto/rand"
	"time"

	db "api.north-path.site/utils/dynamodb"

	"errors"
	"fmt"

	errs "api.north-path.site/utils/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/oklog/ulid"
)

const tableName = "posts"

type Post struct {
	PostId      string   `json:"postId" dynamodbav:"postId"`
	Email       string   `json:"email" dynamodbav:"email"`
	Subject     string   `json:"subject" dynamodbav:"subject"`
	Content     string   `json:"content" dynamodbav:"content"`
	Categories  []string `json:"categories" dynamodbav:"categories"`
	Images      []string `json:"images" dynamodbav:"images"`
	CreatedDate string   `json:"createdDate" dynamodbav:"createdDate"`
	UpdatedDate string   `json:"updatedDate" dynamodbav:"updatedDate"`
	client      *db.Client
}

type PostMethod interface {
	CreateNew(email, subject, content *string, images, categories *[]string) (Post, error)
	FindById(id string) (*Post, error)
	DeleteById(id string) error
	Search(limit int32, currentToken string, category string) ([]Post, *string, error)
}

func New() PostMethod {
	client := db.New(tableName)

	return Post{client: &client}
}

func (p Post) CreateNew(email, subject, content *string, images, categories *[]string) (Post, error) {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	post := &Post{
		PostId:      id.String(),
		Email:       *email,
		Subject:     *subject,
		Content:     *content,
		Categories:  *categories,
		Images:      *images,
		CreatedDate: time.Now().Format(time.RFC3339),
		UpdatedDate: time.Now().Format(time.RFC3339),
	}

	return *post, p.client.CreateOrUpdate(post)
}

func (p Post) FindById(id string) (*Post, error) {
	item, err := p.client.FindById("postId", id)
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(item, &p)
	if err != nil {
		return nil, errors.New(errs.UnmarshalError)
	}

	return &p, nil
}

func (p Post) DeleteById(id string) error {
	return p.client.DeleteById("postId", id)
}

func (p Post) Search(limit int32, currentToken string, category string) ([]Post, *string, error) {

	statement := fmt.Sprintf("SELECT * FROM \"%v\"", *p.client.TableName)
	if category != "" {
		statement = fmt.Sprintf("SELECT * FROM \"%v\" WHERE CONTAINS (categories, '?')", *p.client.TableName)
	}

	input := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(statement),
		Limit:     &limit,
	}

	if currentToken != "" {
		input.NextToken = &currentToken
	}

	if category != "" {
		params, err := attributevalue.MarshalList([]interface{}{category})
		if err != nil {
			return nil, nil, err
		}

		input.Parameters = params
	}

	response, err := p.client.ExecuteStatement(input)

	if err != nil {
		return nil, nil, err
	}

	var posts []Post
	err = attributevalue.UnmarshalListOfMaps(response.Items, &posts)
	nextToken := response.NextToken
	if err != nil {
		return nil, nil, err
	}
	return posts, nextToken, nil

}
