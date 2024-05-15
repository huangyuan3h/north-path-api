package main

import (
	"net/http"

	"api.north-path.site/post/db"
	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type ViewPostBody struct {
	Id string `json:"id" validate:"required"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	viewPostReq := &ViewPostBody{
		Id: request.PathParameters["id"],
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(viewPostReq)

	// verify
	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "Id":
			errMessage = errors.NotValidSubject
		}

		return errors.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	db_client := db.New()

	post, err := db_client.FindById(viewPostReq.Id)

	if err != nil {
		return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	return awsHttp.Ok(post, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
