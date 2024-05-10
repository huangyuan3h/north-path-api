package main

import (
	"net/http"

	"encoding/json"

	"api.north-path.site/post/db"
	"api.north-path.site/post/types"
	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type SearchPostBody struct {
	Limit        int32  `json:"limit" validate:"required,max=100"`
	CurrentToken string `json:"current_token"`
	Category     string `json:"category"`
}

type ViewPostResponse struct {
	Results   []types.Post `json:"results"`
	NextToken *string      `json:"next_token"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	var body SearchPostBody
	err := json.Unmarshal([]byte(request.Body), &body)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(body)

	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "Limit":
			errMessage = errors.NotValidEmail
		}

		return errors.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	db_client := db.New()

	posts, nextToken, err := db_client.Search(body.Limit, body.CurrentToken, body.Category)

	if err != nil {
		return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	return awsHttp.Ok(&ViewPostResponse{
		Results:   posts,
		NextToken: nextToken,
	}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
