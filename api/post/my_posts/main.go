package main

import (
	"encoding/json"
	"net/http"

	"api.north-path.site/post/db"
	"api.north-path.site/post/types"
	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	"api.north-path.site/utils/jwt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type MyPostsRequest struct {
	Email     string `json:"email"`
	Limit     int32  `json:"limit"`
	NextToken string `json:"next_token"`
}

type MyPostsResponse struct {
	Results   []types.Post `json:"results"`
	NextToken string       `json:"next_token"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	myPostsReq := &MyPostsRequest{}
	err := json.Unmarshal([]byte(request.Body), &myPostsReq)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	if myPostsReq.Email == "" {
		claim, err := jwt.VerifyRequest(request)
		if err != nil {
			return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
		}
		myPostsReq.Email = claim.Email
	}

	// validate
	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(myPostsReq)

	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "Email":
			errMessage = errors.NotValidEmail
		case "Limit":
			errMessage = errors.NotValidLimit
		}

		return errors.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	db := db.New()

	posts, nextToken, err := db.FindByEmail(myPostsReq.Limit, myPostsReq.NextToken, myPostsReq.Email)
	if err != nil {
		return errors.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
	}

	return awsHttp.Ok(&MyPostsResponse{
		Results:   posts,
		NextToken: *nextToken,
	}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
