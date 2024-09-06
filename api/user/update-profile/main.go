package main

import (
	"net/http"

	"encoding/json"

	"north-path.it-t.xyz/user/db"
	"north-path.it-t.xyz/user/types"
	"north-path.it-t.xyz/utils/errors"
	awsHttp "north-path.it-t.xyz/utils/http"
	"north-path.it-t.xyz/utils/jwt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

type UpdateProfileBody struct {
	Email    string `json:"email"  validate:"required,email"`
	Avatar   string `json:"avatar" `
	UserName string `json:"userName" validate:"required,min=6,max=50"`
	Bio      string `json:"bio"`
}

type UpdateProfileResponse struct {
	Authorization string `json:"Authorization"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	var updateProfileReq UpdateProfileBody
	err := json.Unmarshal([]byte(request.Body), &updateProfileReq)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	errStruct := validate.Struct(updateProfileReq)

	// verify
	if errStruct != nil {
		firstErr := errStruct.(validator.ValidationErrors)[0]
		var errMessage string
		switch t := firstErr.StructField(); t {
		case "Email":
			errMessage = errors.NotValidEmail
		case "UserName":
			errMessage = errors.UseNameInvalid
		}

		return errors.New(errMessage, http.StatusBadRequest).GatewayResponse()
	}

	claim, err := jwt.VerifyRequest(request)

	if err != nil {
		return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	//verify ownership

	if claim.Email != updateProfileReq.Email {
		return errors.New(errors.OwnerNotMatch, http.StatusBadRequest).GatewayResponse()
	}

	// find the original to db

	db_client := db.New()

	err = db_client.CreateNew(&types.User{
		Email:    updateProfileReq.Email,
		UserName: updateProfileReq.UserName,
		Avatar:   updateProfileReq.Avatar,
		Bio:      updateProfileReq.Bio,
	})

	if err != nil {
		return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	jwtObj := map[string]interface{}{
		"email":    updateProfileReq.Email,
		"avatar":   updateProfileReq.Avatar,
		"userName": updateProfileReq.UserName,
	}
	jwt_token, err := jwt.CreateToken(jwtObj)

	if err != nil {
		return errors.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
	}

	return awsHttp.Ok(&UpdateProfileResponse{
		Authorization: jwt_token,
	}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
