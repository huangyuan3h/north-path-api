package main

import (
	"net/http"

	"encoding/json"

	errors "api.north-path.site/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type CreateAccountBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Handler(request events.APIGatewayV2HTTPRequest)(events.APIGatewayProxyResponse, error) {

	var user CreateAccountBody
  	err := json.Unmarshal([]byte(request.Body), &user)

	if err!= nil {
		return events.APIGatewayProxyResponse{
            StatusCode: http.StatusBadRequest,
            Body:       errors.JSONParseError,
        }, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "created",
		StatusCode: http.StatusCreated,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
