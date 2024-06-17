package main

import (
	"encoding/json"
	"net/http"

	"api.north-path.site/message/db"
	"api.north-path.site/message/types"
	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type SendMessageBody struct {
	Subject   string `json:"subject" validate:"required,min=6,max=50"`
	Content   string `json:"content" validate:"required,min=6,max=500"`
	ToEmail   string `json:"toEmail" `
	FromEmail string `json:"fromEmail"  validate:"required,email"`
}

type ContactAdminResponse struct {
	Message string `json:"message"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var sendMessageBody SendMessageBody
	err := json.Unmarshal([]byte(request.Body), &sendMessageBody)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(sendMessageBody)

	// verify
	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "FromEmail":
			errMessage = errors.NotValidEmail
		case "Subject":
			errMessage = errors.SubjectInvalid
		case "Content":
			errMessage = errors.ContentInvalid
		}

		return errors.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	// save the log to db

	client := db.New()
	client.CreateNew(&types.Message{
		FromEmail: sendMessageBody.FromEmail,
		ToEmail:   sendMessageBody.ToEmail,
		Subject:   sendMessageBody.Subject,
		Content:   sendMessageBody.Content,
	})

	return awsHttp.Ok(&ContactAdminResponse{
		Message: "ok",
	}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
