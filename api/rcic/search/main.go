package main

import (
	"encoding/json"
	"net/http"

	"api.north-path.site/utils/errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RCICRequest struct {
	RCIC string `json:"rcic" validate:"required"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	var rcicReq RCICRequest
	err := json.Unmarshal([]byte(request.Body), &rcicReq)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	return events.APIGatewayProxyResponse{
		Body:       "ok",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
