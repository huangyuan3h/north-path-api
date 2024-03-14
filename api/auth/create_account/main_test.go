package main

import (
	"testing"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

  func TestHandlerSanity(t *testing.T) {
	input := events.APIGatewayV2HTTPRequest{
		Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: "{\"name\":\"John\"}",
	}

	result, _ := Handler(input)

	if result.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201, got %d", result.StatusCode)
	}

  }

  func TestHandlerJSONParse(t *testing.T) {
	input := events.APIGatewayV2HTTPRequest{
		Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: "not a json",
	}

	result, _ := Handler(input)

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", result.StatusCode)
	}

	if result.Body != "JSON Parse Error" {
		t.Errorf("Expected JSON Parse Error, got %s", result.Body)
	}

  }