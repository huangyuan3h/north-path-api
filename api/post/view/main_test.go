package main

import (
	"testing"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerSanity(t *testing.T) {

	input := events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"id": "01HW4V8ZGWWT1SZ1VXJ94TACYE",
		},
	}

	result, _ := Handler(input)

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", result.StatusCode)
	}
}
