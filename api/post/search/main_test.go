package main

import (
	"testing"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerSanity(t *testing.T) {

	input := events.APIGatewayV2HTTPRequest{
		Body: "{\"limit\":1,\"next_token\":\"\",\"category\":\"asd\"}",
	}

	result, _ := Handler(input)

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", result.StatusCode)
	}
}
