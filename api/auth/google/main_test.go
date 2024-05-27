package main

import (
	"testing"

	"net/http"

	"os"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerSanity(t *testing.T) {

	os.Setenv("JWT_SECRET", "h3OOumyH3vLgUhHve7bLPv8hgNXbxUQr")
	os.Setenv("AVATAR_BUCKET_NAME", "dev-north-path-api-stack-avatarbucketd80dbdb5-hctiiwefencj")
	os.Setenv("GOOGLE_CLIENT_ID", "370355018861-c6ukhtuk0c9tstbi7k3ir5buhnk9lmvr.apps.googleusercontent.com")
	credential := "eyJhbGciOiJSUzI1NiIsImtpZCI6IjMyM2IyMTRhZTY5NzVhMGYwMzRlYTc3MzU0ZGMwYzI1ZDAzNjQyZGMiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiIzNzAzNTUwMTg4NjEtYzZ1a2h0dWswYzl0c3RiaTdrM2lyNWJ1aG5rOWxtdnIuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiIzNzAzNTUwMTg4NjEtYzZ1a2h0dWswYzl0c3RiaTdrM2lyNWJ1aG5rOWxtdnIuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDMwNzI5OTQ3ODU2MTQzMjc3MDEiLCJlbWFpbCI6Imh1YW5neXVhbjNoQGdtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYmYiOjE3MTYzNjk2MTQsIm5hbWUiOiLpu4TnvJjvvIhZdWFuIEh1YW5n77yJIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FDZzhvY0xCa0RPRXY5N3Q0bUpCdUJRT3pSWmtFSktsWEtTazdzOWN1dVZxTUwxd2Z0MGdIZlphPXM5Ni1jIiwiZ2l2ZW5fbmFtZSI6Iue8mCIsImZhbWlseV9uYW1lIjoi6buEIiwiaWF0IjoxNzE2MzY5OTE0LCJleHAiOjE3MTYzNzM1MTQsImp0aSI6IjFlMTgzYjliNjAyYzQ2OWU0YjlkNTM3NzhmODAzOTk4MDc0ZGJlNWYifQ.mSGGxKtDm8lQQrZyZ6QeX9HeXxLkzYTQiHOYmWdONASxvq_tp495DzllFoeh_n4WPSObnJSqiQpqX9PKjcRI3Dkb1A9ADaL1aPgGe_c-a_NlfEaBEh755LC9HZ3ZqcXIgE-2e6nsrVbcF9OWgs8Ck19RO91Xiwdb4tknEAPevmUGsKKcqxYQHVyP6CmCfoHhpp53FKCLwnklFTPZ_XrzVBQ13g9FLKulY5sLEnYxBm1Hc3Vnk_EtxGqGEUl3XnhU-WPtedgnJsOiscf055HRdHI8FnW_a7M82yJGF8E9m8v1BS2dCJ4vgv5Pm6pISDo6zYUfsZ50_0cpYOPhgXmIKw"
	input := events.APIGatewayProxyRequest{
		Body: "{\"credential\":\"" + credential + "\"}",
	}

	result, _ := handleGoogleLogin(input)
	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", result.StatusCode)
	}
}