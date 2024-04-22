package main

import (
	"strings"

	"errors"

	"api.north-path.site/utils/jwt"
	"github.com/aws/aws-lambda-go/events"
)

func VerifyAuth(request events.APIGatewayV2HTTPRequest) (*jwt.MyClaims, error) {

	authStr := request.Headers["Authorization"]

	if authStr == "" {

		for _, cookie := range request.Cookies {

			parts := strings.Split(cookie, "=")

			if parts[0] == "Authorization" {
				authStr = parts[1]
				break
			}
		}
	}

	if authStr == "" {
		return nil, errors.New("authorization is not found in request")
	}

	myClaims, err := jwt.VerifyToken(authStr)

	if err != nil {
		return nil, err
	}

	return myClaims, nil
}
