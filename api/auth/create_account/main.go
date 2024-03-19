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
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(acocuntReq)

	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "Email":
			errMessage = errors.NotValidEmail
		case "Password":
			errMessage = errors.PasswordError
		}

		return errors.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	// detail validation
	var regContainsLow = regexp.MustCompile("[a-z]+")
	var regContainsUpper = regexp.MustCompile("[A-Z]+")
	var regContainsNumber = regexp.MustCompile("[0-9]+")

	if !regContainsLow.MatchString(acocuntReq.Password) || !regContainsUpper.MatchString(acocuntReq.Password) || !regContainsNumber.MatchString(acocuntReq.Password) {
		return errors.New(errors.PasswordError, http.StatusBadRequest).GatewayResponse()
	}

	auth := auth.New()

	// add record and send email

	err = auth.CreateAccount(&acocuntReq.Email, &acocuntReq.Password)

	if err != nil {
		return errors.New(errors.InsertDBError, http.StatusBadRequest).GatewayResponse()
	}

	return events.APIGatewayProxyResponse{
		Body:       "created",
		StatusCode: http.StatusCreated,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
