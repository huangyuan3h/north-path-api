package main

import (
	"net/http"

	"encoding/json"
	"regexp"

	"api.north-path.site/auth/db/auth"
	user "api.north-path.site/user/db"
	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type CreateAccountBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type CreateAccountResponse struct {
	Message string `json:"message"`
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
	user := user.New()

	// add record and send email

	err = auth.CreateAccount(&acocuntReq.Email, &acocuntReq.Password)

	if err != nil {
		return errors.New(errors.InsertDBError, http.StatusBadRequest).GatewayResponse()
	}

	err = user.CreateNew(&acocuntReq.Email)
	if err != nil {
		return errors.New(errors.InsertDBError, http.StatusBadRequest).GatewayResponse()
	}

	return awsHttp.Ok(CreateAccountResponse{Message: "created"}, http.StatusCreated)
}

func main() {
	lambda.Start(Handler)
}
