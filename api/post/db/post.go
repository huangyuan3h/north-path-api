package db

import (
	"crypto/rand"

	"time"

	db "north-path.it-t.xyz/utils/dynamodb"

	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/oklog/ulid"
	types "north-path.it-t.xyz/post/types"
)

const tableName = "posts"

type Post struct {
	types.Post
	client *db.Client
}

type PostMethod interface {
	CreateOrUpdate(pid, email, subject, content, category, location, bilibili, youtube *string, images, topics *[]string) (*types.Post, error)
	FindById(id string) (*types.Post, error)
	DeleteById(id string) error
	Search(limit int32, currentId string, category string) ([]types.Post, *string, error)
	FindByEmail(limit int32, currentId string, email string) ([]types.Post, *string, error)
}

func New() PostMethod {
	client := db.New(tableName)

	return Post{client: &client}
}

func (p Post) CreateOrUpdate(pid, email, subject, content, category, location, bilibili, youtube *string, images, topics *[]string) (*types.Post, error) {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	postId := id.String()
	if pid != nil && *pid != "" {
		postId = *pid
	}

	timestamp_ms := int64(time.Now().UnixNano() / int64(time.Millisecond))
	post :=
		&types.Post{
			PostId:       postId,
			Email:        *email,
			Subject:      *subject,
			Content:      *content,
			Topics:       *topics,
			Category:     *category,
			Location:     *location,
			Images:       *images,
			Bilibili:     *bilibili,
			Youtube:      *youtube,
			UpdatedDate:  time.Now().Format(time.RFC3339),
			Like:         0,
			SortingScore: timestamp_ms,
			Status:       "Active",
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
		statement = fmt.Sprintf("SELECT * FROM \"%v\".\"all\" where status = 'Active' order by sortingScore desc", *p.client.TableName)
	} else {
		statement = fmt.Sprintf("SELECT * FROM \"%v\".\"category\" where category = ? order by sortingScore desc", *p.client.TableName)
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

func (p Post) FindByEmail(limit int32, currentId string, email string) ([]types.Post, *string, error) {
	statement := fmt.Sprintf("SELECT * FROM \"%v\".\"myPost\" where email = ? order by updatedDate desc", *p.client.TableName)

	emptyStr := ""

	input := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(statement),
		Limit:     aws.Int32(limit),
	}

	if currentId != "" {
		input.NextToken = aws.String(currentId)
	}

	params, err := attributevalue.MarshalList([]interface{}{email})

	if err != nil {
		return nil, &emptyStr, err
	}
	input.Parameters = params
	response, err := p.client.ExecuteStatement(input)

	if err != nil {
		return nil, &emptyStr, err
	}

	var posts []types.Post
	err = attributevalue.UnmarshalListOfMaps(response.Items, &posts)
	nextToken := response.NextToken
	if err != nil {
		return nil, &emptyStr, err
	}

	if response.NextToken == nil {
		return posts, &emptyStr, nil
	}
	return posts, nextToken, nil
}
