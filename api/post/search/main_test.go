package main

import (
	"testing"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerSanity(t *testing.T) {

	input := events.APIGatewayV2HTTPRequest{
		Body: "{\"limit\":10,\"next_token\":\"\",\"category\":\"\"}",
	}

	result, _ := Handler(input)

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", result.StatusCode)
	}
}
