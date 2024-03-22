package main

import (
	"encoding/json"
	"net/http"

	"api.north-path.site/auth/db/auth"
	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	"api.north-path.site/utils/jwt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type LoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type LoginResponse struct {
	Authorization string `json:"Authorization"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	var loginReq LoginBody
	err := json.Unmarshal([]byte(request.Body), &loginReq)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(loginReq)

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

	auth := auth.New()

	// add record and send email

	err = auth.VerifyLogin(&loginReq.Email, &loginReq.Password)

	if err != nil {
		return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	// generate jwt token

	jwtObj := map[string]interface{}{
		"email": loginReq.Email,
	}
	jwt_token, err := jwt.CreateToken(jwtObj)

	if err != nil {
		return errors.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
	}

	return awsHttp.Ok(LoginResponse{Authorization: jwt_token}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
