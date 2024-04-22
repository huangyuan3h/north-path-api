package main

import (
	"net/http"

	"encoding/json"

	"api.north-path.site/post/db"
	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	"api.north-path.site/utils/jwt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type CreatePostBody struct {
	Subject    string   `json:"subject" validate:"required,min=6,max=50"`
	Content    string   `json:"content"  validate:"required,max=5000"`
	Categories []string `json:"categories" validate:"required"`
	Images     []string `json:"images"  validate:"required"`
}

type CreatePostResponse struct {
	Id         string   `json:"id"`
	Subject    string   `json:"subject" `
	Content    string   `json:"content" `
	Categories []string `json:"categories"`
	Images     []string `json:"images"`
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

		return errors.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	claim, err := jwt.VerifyRequest(request)

	if err != nil {
		return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	// save to db

	db_client := db.New()

	post, err := db_client.CreateNew(&claim.Email, &createPostReq.Subject, &createPostReq.Content, &createPostReq.Images, &createPostReq.Categories)

	if err != nil {
		return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	return awsHttp.Ok(&CreatePostResponse{
		Id:         post.PostId,
		Subject:    createPostReq.Subject,
		Content:    createPostReq.Content,
		Categories: createPostReq.Categories,
		Images:     createPostReq.Images,
	}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
