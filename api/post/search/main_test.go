package main

import (
	"testing"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerSanity(t *testing.T) {

	input := events.APIGatewayV2HTTPRequest{
		QueryStringParameters: map[string]string{
			"limit":         "1",
			"current_token": "/n5Y8KGG6kBkkvJ/QTVQo74ZbTPWBOm/N93u6UkDeeljz2YdD/OqNqcFs+DRBiMXsrS3/UNiGfG+UytGo3UR9ug+cGpe2ontdGzVDDMVlfml3FfJf5xgOxFAY3FPDikOQcbIqdtqtsesUeoDuKWT0CWtLnt+E14W6jA7YPq5Oydl2h9uMXAy08YU/GpCo38TdEyy3hQMrPe+Ru8Re9Lp4GnGBuec8L/8yoKu4Jz99vhOeM/p3T0hU7mMOU/hVrA4uxdAXfTiv7YVEY077wIBCf7yKVW7VKlpZNvxWIb7qSh5tc14quVTToYa1YiendK3pcx5p34SHANq5U6du/7nrdOk9AJT2yqK0ScydX2mvu9WqEzr7Dj+4Fz9jqzK1Q==",
		},
	}

	result, _ := Handler(input)

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", result.StatusCode)
	}
}
