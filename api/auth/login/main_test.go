package main

import (
	"testing"

	"net/http"

	"os"

	errors "north-path.it-t.xyz/utils/errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-cmp/cmp"
)

func TestHandlerSanity(t *testing.T) {
	os.Setenv("AUTH_SECRET", "GLbR3zUjXPbSKLwsSqNDTG3ODNkZYDdF")
	input := events.APIGatewayV2HTTPRequest{
		Body: "{\"email\":\"abc123@qq.com\", \"password\":\"Password123\"}",
	}

	result, _ := Handler(input)

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", result.StatusCode)
	}
}

func TestHandlerJSONParse(t *testing.T) {
	input := events.APIGatewayV2HTTPRequest{
		Body: "not a json",
	}

	result, _ := Handler(input)

	httpError, err := errors.StringToErrorMessage(result.Body)

	if err != nil {
		t.Errorf("response is not a JSON %d", err)
	}

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", result.StatusCode)
	}

	if !cmp.Equal(httpError.Message, errors.JSONParseError) {
		t.Errorf("Expected %s, got %s", errors.JSONParseError, httpError.Message)
	}
}

func TestHandlerNotValidEmail(t *testing.T) {
	input := events.APIGatewayV2HTTPRequest{
		Body: "{\"email\":\"not a email\", \"password\":\"P@ssw0rd!\"}",
	}

	result, _ := Handler(input)

	httpError, err := errors.StringToErrorMessage(result.Body)

	if err != nil {
		t.Errorf("response is not a JSON %d", err)
	}

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", result.StatusCode)
	}

	if !cmp.Equal(httpError.Message, errors.NotValidEmail) {
		t.Errorf("Expected %s, got %s", errors.NotValidEmail, httpError.Message)
	}
}

func TestHandlerNotValidPassword(t *testing.T) {

	// password length less than 6

	passwords := [1]string{"error"}

	for _, password := range passwords {

		input := events.APIGatewayV2HTTPRequest{
			Body: "{\"email\":\"abc123@qq.com\", \"password\":\"" + password + "\"}",
		}

		result, _ := Handler(input)

		httpError, err := errors.StringToErrorMessage(result.Body)

		if err != nil {
			t.Errorf("response is not a JSON %d", err)
		}

		if result.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", result.StatusCode)
		}

		if !cmp.Equal(httpError.Message, errors.PasswordError) {
			t.Errorf("Expected %s, got %s", errors.PasswordError, httpError.Message)
		}
	}
}
