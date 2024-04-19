package main

import (
	"net/http"

	"encoding/json"

	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	// "api.north-path.site/post/db"
)

type CreatePostBody struct {
	Subject    string   `json:"subject" validate:"required,min=6,max=50"`
	Content    string   `json:"content"  validate:"required,max=5000"`
	Categories []string `json:"categories" validate:"required"`
	Images     []string `json:"images"  validate:"required"`
}

type CreatePostResponse struct {
	Subject    string   `json:"subject" validate:"required,min=6,max=50"`
	Content    string   `json:"content"  validate:"required,max=5000"`
	Categories []string `json:"categories" validate:"required"`
	Images     []string `json:"images"  validate:"required"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	var createPostReq CreatePostBody
	err := json.Unmarshal([]byte(request.Body), &createPostReq)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(createPostReq)

	// verify
	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "Subject":
			errMessage = errors.NotValidSubject
		case "Content":
			errMessage = errors.NotValidContent
		case "Categories":
			errMessage = errors.NotValidCategories
		case "Images":
			errMessage = errors.NotValidImages
		}

		// save to db

		// db_client := db.New()

		// db_client.CreateNew()

		return errors.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	return awsHttp.Ok(&CreatePostResponse{
		Subject:    createPostReq.Subject,
		Content:    createPostReq.Content,
		Categories: createPostReq.Categories,
		Images:     createPostReq.Images,
	}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
