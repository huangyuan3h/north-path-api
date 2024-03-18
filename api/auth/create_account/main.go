package main

import (
	"net/http"

	"encoding/json"
	"regexp"

	"api.north-path.site/auth/db/auth"
	"api.north-path.site/utils/errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type CreateAccountBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	var acocuntReq CreateAccountBody
	err := json.Unmarshal([]byte(request.Body), &acocuntReq)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       errors.JSONParseError,
		}, nil
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(acocuntReq)

	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var body string
		switch t := firstErr.StructField(); t {
		case "Email":
			body = errors.NotValidEmail
		case "Password":
			body = errors.PasswordError
		}

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       body,
		}, nil

	}

	// detail validation
	var regContainsLow = regexp.MustCompile("[a-z]+")
	var regContainsUpper = regexp.MustCompile("[A-Z]+")
	var regContainsNumber = regexp.MustCompile("[0-9]+")

	if !regContainsLow.MatchString(acocuntReq.Password) || !regContainsUpper.MatchString(acocuntReq.Password) || !regContainsNumber.MatchString(acocuntReq.Password) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       errors.PasswordError,
		}, nil
	}

	auth := auth.New()

	// check if the email is existed

	// add record and send email

	err = auth.CreateAccount(&acocuntReq.Email, &acocuntReq.Password)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       errors.InsertDBError,
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
