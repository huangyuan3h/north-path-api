package main

import (
	"context"
	"encoding/json"
	sysErrors "errors"
	"net/http"
	"os"
	"strings"

	"north-path.it-t.xyz/post/db"
	"north-path.it-t.xyz/utils/errors"
	awsHttp "north-path.it-t.xyz/utils/http"
	"north-path.it-t.xyz/utils/jwt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type DeletePostRequest struct {
	PostIds []string `json:"post_ids"`
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

// todo: optimize to batch operations

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

	for _, pid := range deleteReq.PostIds {

		post, err := db.FindById(pid)

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

		err = db.DeleteById(pid)
		if err != nil {
			return errors.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
		}
	}

	return awsHttp.Ok(&DeletePostResponse{
		Message: "success",
	}, http.StatusOK)
}

func main() {
	lambda.Start(Handler)
}
