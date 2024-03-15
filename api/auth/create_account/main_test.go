package main

import (
	"testing"

	"net/http"

	errors "api.north-path.site/utils/errors"
	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerSanity(t *testing.T) {
	input := events.APIGatewayV2HTTPRequest{
		Body: "{\"email\":\"abc123@qq.com\", \"password\":\"Password123\"}",
	}

	result, _ := Handler(input)

	if result.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201, got %d", result.StatusCode)
	}
}

func TestHandlerJSONParse(t *testing.T) {
	input := events.APIGatewayV2HTTPRequest{
		Body: "not a json",
	}

	result, _ := Handler(input)

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", result.StatusCode)
	}

	if result.Body != errors.JSONParseError {
		t.Errorf("Expected %s, got %s", errors.JSONParseError, result.Body)
	}
}

func TestHandlerNotValidEmail(t *testing.T) {
	input := events.APIGatewayV2HTTPRequest{
		Body: "{\"email\":\"not a email\", \"password\":\"P@ssw0rd!\"}",
	}

	result, _ := Handler(input)

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", result.StatusCode)
	}

	if result.Body != errors.NotValidEmail {
		t.Errorf("Expected %s, got %s", errors.NotValidEmail, result.Body)
	}
}

func TestHandlerNotValidPassword(t *testing.T) {

	// password length less than 6
	//password no upcase
	//password no lowcase
	//password no number

	passwords := [4]string{"error", "password123", "PASSWORD123", "PASSWORD"}

	for _, password := range passwords {

		input := events.APIGatewayV2HTTPRequest{
			Body: "{\"email\":\"abc123@qq.com\", \"password\":\"" + password + "\"}",
		}

		result, _ := Handler(input)

		if result.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", result.StatusCode)
		}

		if result.Body != errors.PasswordError {
			t.Errorf("Expected %s, got %s", errors.PasswordError, result.Body)
		}
	}
}
