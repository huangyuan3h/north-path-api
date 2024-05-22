package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"context"
	"encoding/base64"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"path/filepath"

	"os"

	user "api.north-path.site/user/db"
	"api.north-path.site/utils/errors"
	awsHttp "api.north-path.site/utils/http"
	myJWT "api.north-path.site/utils/jwt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type GoogleLoginRequest struct {
	Credential string `json:"credential"`
}

type LoginResponse struct {
	Authorization string `json:"Authorization"`
}

type GoogleUserType struct {
	Email   string  `json:"email"`
	Picture string  `json:"picture"`
	Name    string  `json:"name"`
	Exp     float64 `json:"exp"`
}

func ExtractUserInfo(token string) (GoogleUserType, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return GoogleUserType{}, fmt.Errorf("invalid token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return GoogleUserType{}, fmt.Errorf("error decoding payload: %w", err)
	}

	var userInfo GoogleUserType
	err = json.Unmarshal(payload, &userInfo)
	if err != nil {
		return GoogleUserType{}, fmt.Errorf("error unmarshalling user info: %w", err)
	}

	exp := userInfo.Exp

	if time.Now().Unix() > int64(exp) {
		return GoogleUserType{}, fmt.Errorf("token expired")
	}

	return userInfo, nil
}

func uploadToS3(ctx context.Context, url string, bucket string, key string) error {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := s3.NewFromConfig(cfg)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read image data: %w", err)
	}

	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to decode image config: %w", err)
	}
	if format != "jpeg" && format != "png" {
		return fmt.Errorf("unsupported image format: %s", format)
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentType:   aws.String(fmt.Sprintf("image/%s", format)),
		ContentLength: aws.Int64(int64(len(data))),
		ACL:           types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return fmt.Errorf("failed to upload image to S3: %w", err)
	}

	return nil
}

func handleGoogleLogin(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var loginReq GoogleLoginRequest
	err := json.Unmarshal([]byte(request.Body), &loginReq)
	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	userInfo, err := ExtractUserInfo(loginReq.Credential)
	if err != nil {
		return errors.New(errors.JSONParseError, http.StatusBadRequest).GatewayResponse()
	}

	userRepo := user.New()
	u, err := userRepo.FindByEmail(&userInfo.Email)
	if err != nil {
		return errors.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
	}

	if u.Email == "" { // not found
		bucketName := os.Getenv("NEXT_PUBLIC_BUCKET_NAME")
		fileName := fmt.Sprintf("%s-%d%s", uuid.New(), time.Now().Unix(), filepath.Ext(userInfo.Picture))

		err = uploadToS3(ctx, userInfo.Picture, bucketName, fileName)
		if err != nil {
			fmt.Println("Failed to upload image to S3:", err)
		}

		if err != nil {
			return errors.New(fmt.Sprintf("failed to upload image to S3: %s", err.Error()), http.StatusInternalServerError).GatewayResponse()
		}

		u = &user.User{
			Email:    userInfo.Email,
			UserName: userInfo.Name,
			Avatar:   fmt.Sprintf("https://%s.s3.us-east-1.amazonaws.com/%s", bucketName, fileName),
		}
		if err := userRepo.CreateNew(u); err != nil {
			return errors.New("Failed to create user", http.StatusInternalServerError).GatewayResponse()
		}
	}

	jwtObj := map[string]interface{}{
		"email":    u.Email,
		"avatar":   u.Avatar,
		"userName": u.UserName,
	}
	jwtToken, err := myJWT.CreateToken(jwtObj)
	if err != nil {
		return errors.New(err.Error(), http.StatusInternalServerError).GatewayResponse()
	}

	return awsHttp.Ok(LoginResponse{Authorization: jwtToken}, http.StatusOK)
}

func main() {
	lambda.Start(handleGoogleLogin)
}
