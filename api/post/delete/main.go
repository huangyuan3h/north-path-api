package main

import (
	"context"
	"encoding/json"
	sysErrors "errors"
	"net/http"
	"os"
	"strings"

	"api.north-path.site/post/db"
	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	"api.north-path.site/utils/jwt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DeletePostRequest struct {
	PostId string `json:"post_id"`
}

type DeletePostResponse struct {
	Message string `json:"message"`
}

func extractImageKey(imageUrl string) (string, error) {
	if !strings.HasPrefix(imageUrl, "https://") {
		return "", sysErrors.New(errors.InvalidImageUrl)
	}

	imageUrl = imageUrl[8:]

	parts := strings.Split(imageUrl, "/")

	if len(parts) != 2 {
		return "", sysErrors.New(errors.InvalidImageUrl)
	}

	key := parts[1]

	return key, nil
}

func deleteImage(images []string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	bucket := os.Getenv("POST_IMAGE_BUCKET_NAME")
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)
	for _, imageUrl := range images {

		key, err := extractImageKey(imageUrl)
		if err != nil {

			return err
		}

		input := &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}

		_, err = client.DeleteObject(context.TODO(), input)
		if err != nil {
			return err
		}
	}

	return nil
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {

	deleteReq := &DeletePostRequest{}
	err := json.Unmarshal([]byte(request.Body), &deleteReq)

	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	claim, err := jwt.VerifyRequest(request)
	if err != nil {
		return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	email := claim.Email

	db := db.New()

	post, err := db.FindById(deleteReq.PostId)

	if err != nil {
		return errors.New(err.Error(), http.StatusBadRequest).GatewayResponse()
	}

	// make sure the owner is the current user
	if post.Email != email {
		return errors.New(errors.OwnerNotMatch, http.StatusBadRequest).GatewayResponse()
	}

	// delete the image from s3
	err = deleteImage(post.Images)
	if err != nil {
		return errors.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
	}

	err = db.DeleteById(deleteReq.PostId)

	if err != nil {
		return errors.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
	}

	return awsHttp.Ok(&DeletePostResponse{
		Message: "success",
	}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
