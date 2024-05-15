package db

import (
	"crypto/rand"

	"time"

	db "api.north-path.site/utils/dynamodb"

	"fmt"

	types "api.north-path.site/post/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/oklog/ulid"
)

const tableName = "posts"

type Post struct {
	types.Post
	client *db.Client
}

type PostMethod interface {
	CreateNew(email, subject, content *string, images, categories *[]string) (*types.Post, error)
	FindById(id string) (*types.Post, error)
	DeleteById(id string) error
	Search(limit int32, currentId string, category string) ([]types.Post, *string, error)
}

func New() PostMethod {
	client := db.New(tableName)

	return Post{client: &client}
}

func (p Post) CreateNew(email, subject, content *string, images, categories *[]string) (*types.Post, error) {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	post :=
		&types.Post{
			PostId:      id.String(),
			Email:       *email,
			Subject:     *subject,
			Content:     *content,
			Categories:  *categories,
			Images:      *images,
			CreatedDate: time.Now().Format(time.RFC3339),
			UpdatedDate: time.Now().Format(time.RFC3339),
			Status:      "Active",
		}

	return post, p.client.CreateOrUpdate(post)
}

func (p Post) FindById(id string) (*types.Post, error) {
	statement := fmt.Sprintf("SELECT * FROM \"%v\" WHERE postId=?", *p.client.TableName)

	params, err := attributevalue.MarshalList([]interface{}{id})

	if err != nil {
		return nil, err
	}

	input := &dynamodb.ExecuteStatementInput{
		Statement:  aws.String(statement),
		Parameters: params,
	}
	response, err := p.client.ExecuteStatement(input)

	if err != nil {
		return nil, err
	}

	var post types.Post
	err = attributevalue.UnmarshalMap(response.Items[0], &post)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p Post) DeleteById(id string) error {
	return p.client.DeleteById("postId", id)
}

func (p Post) Search(limit int32, currentId string, category string) ([]types.Post, *string, error) {

	var statement string
	if category == "" {
		statement = fmt.Sprintf("SELECT * FROM \"%v\".\"GSI1\" where status = 'Active' order by updatedDate desc", *p.client.TableName)
	} else {
		statement = fmt.Sprintf("SELECT * FROM \"%v\".\"GSI1\" where status = 'Active' and contains(\"categories\", ?) order by updatedDate desc", *p.client.TableName)
	}

	input := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(statement),
		Limit:     aws.Int32(limit),
	}

	if currentId != "" {
		input.NextToken = aws.String(currentId)
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

	var posts []types.Post
	err = attributevalue.UnmarshalListOfMaps(response.Items, &posts)
	nextToken := response.NextToken
	if err != nil {
		return nil, nil, err
	}

	return posts, nextToken, nil
}
