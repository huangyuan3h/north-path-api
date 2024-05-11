package db

import (
	"crypto/rand"

	"time"

	"context"

	db "api.north-path.site/utils/dynamodb"

	"fmt"

	types "api.north-path.site/post/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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
	Search(limit int32, currentId string, category string) ([]types.Post, error)
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

func (p Post) Search(limit int32, currentId string, category string) ([]types.Post, error) {

	var posts []types.Post
	var err error

	projEx := expression.NamesList(
		expression.Name("postId"), expression.Name("email"), expression.Name("subject"), expression.Name("content"))

	expr, err := expression.NewBuilder().WithProjection(projEx).Build()

	if err != nil {
		return nil, err
	}

	scanInput := &dynamodb.ScanInput{
		TableName:                 aws.String(*p.client.TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		Limit:                     aws.Int32(limit),
	}

	if currentId != "" {
		previousPost, err := p.FindById(currentId)
		if err != nil {
			return nil, err
		}
		pointer := types.SearchKeys{PostId: previousPost.PostId, UpdatedDate: previousPost.UpdatedDate}
		startKey, err := attributevalue.MarshalMap(pointer)
		if err != nil {
			return nil, err
		}
		scanInput.ExclusiveStartKey = startKey
	}

	paginator := dynamodb.NewScanPaginator(p.client.Client, scanInput)

	count := 0
	for paginator.HasMorePages() {
		response, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		var postsPage []types.Post
		err = attributevalue.UnmarshalListOfMaps(response.Items, &postsPage)
		if err != nil {
			return nil, err
		}

		count = count + len(postsPage)

		posts = append(posts, postsPage...)
		if count >= int(limit) {
			break
		}
	}

	posts = posts[:limit]

	return posts, nil
}
