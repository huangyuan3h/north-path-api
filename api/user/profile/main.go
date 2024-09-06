package main

import (
	"net/http"

	"north-path.it-t.xyz/user/db"
	"north-path.it-t.xyz/user/types"
	"north-path.it-t.xyz/utils/errors"
	awsHttp "north-path.it-t.xyz/utils/http"
	"north-path.it-t.xyz/utils/jwt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ProfileResponse struct {
	types.User
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	email := request.QueryStringParameters["email"]

	if email == "" {
		claim, err := jwt.VerifyRequest(request)
		if err != nil {
			return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
		}
		email = claim.Email
	}

	db := db.New()

	profile, err := db.FindByEmail(&email)
	if err != nil {
		return errors.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
	}

	return awsHttp.Ok(profile, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
