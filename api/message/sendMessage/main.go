package main

import (
	"context"
	"encoding/json"
	"net/http"

	"api.north-path.site/message/db"
	"api.north-path.site/message/types"
	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	t "github.com/aws/aws-sdk-go-v2/service/ses/types"
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

const DEFAULT_EMAIL = "huangyuan3h@gmail.com"

func sendEmailWithSES(data SendMessageBody) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	client := ses.NewFromConfig(cfg)

	input := &ses.SendEmailInput{
		Destination: &t.Destination{
			ToAddresses: []string{
				data.ToEmail,
			},
		},
		Message: &t.Message{
			Body: &t.Body{
				Html: &t.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(data.Content),
				},
			},
			Subject: &t.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(data.Subject),
			},
		},
		Source: aws.String(DEFAULT_EMAIL),
	}

	_, err = client.SendEmail(context.TODO(), input)
	return err
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var sendMessageBody SendMessageBody
	err := json.Unmarshal([]byte(request.Body), &sendMessageBody)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}
	if sendMessageBody.ToEmail == "" {
		sendMessageBody.ToEmail = DEFAULT_EMAIL
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

	sendMessageBody.Content = sendMessageBody.Content + "<br/> <br/> 消息来自：" + sendMessageBody.FromEmail

	err = sendEmailWithSES(sendMessageBody)

	if err != nil {
		return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	return awsHttp.Ok(&ContactAdminResponse{
		Message: "ok",
	}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
