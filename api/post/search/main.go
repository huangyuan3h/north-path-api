package main

import (
	"net/http"

	"strconv"

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

func atoi32(val string) (int32, error) {
	i, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	body := &SearchPostBody{
		Limit:        10,
		CurrentToken: "",
		Category:     "",
	} // defualt value

	for key, value := range request.QueryStringParameters {
		switch key {
		case "limit":
			limit, err := atoi32(value)
			if err != nil {
				return errors.New("error parsing limit", http.StatusBadRequest).GatewayResponse()
			}
			body.Limit = int32(limit)
		case "current_token":
			body.CurrentToken = value
		case "category":
			body.Category = value
		}
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
