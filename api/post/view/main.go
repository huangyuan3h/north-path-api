package main

import (
	"net/http"

	"encoding/json"

	"api.north-path.site/post/db"
	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type ViewPostBody struct {
	Id string `json:"id" validate:"required"`
}

type ViewPostResponse struct {
	Id         string   `json:"id"`
	Subject    string   `json:"subject" `
	Content    string   `json:"content" `
	Categories []string `json:"categories"`
	Images     []string `json:"images"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	var viewPostReq ViewPostBody
	err := json.Unmarshal([]byte(request.Body), &viewPostReq)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(viewPostReq)

	// verify
	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "Id":
			errMessage = errors.NotValidSubject
		}

		return errors.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	db_client := db.New()

	post, err := db_client.FindById(viewPostReq.Id)

	if err != nil {
		return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	return awsHttp.Ok(&ViewPostResponse{
		Id:         post.PostId,
		Subject:    post.Subject,
		Content:    post.Content,
		Categories: post.Categories,
		Images:     post.Images,
	}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
