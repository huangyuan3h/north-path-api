package main

import (
	"os"
	"testing"

	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandlerSanity(t *testing.T) {

	os.Setenv("JWT_SECRET", "h3OOumyH3vLgUhHve7bLPv8hgNXbxUQr")
	input := events.APIGatewayV2HTTPRequest{
		Body: "{\"limit\":10,\"next_token\":\"\"}",
		Cookies: []string{
			"Authorization=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdmF0YXIiOiJodHRwczovL2Rldi1ub3J0aC1wYXRoLWFwaS1zdGFjay1hdmF0YXJidWNrZXRkODBkYmRiNS1oY3RpaXdlZmVuY2ouczMudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vY2I0NjhhOGEtODliMi00ZjcwLTgzMmQtZmMwODljNWMzZmI5LTE3MTYzNzMwNTEiLCJlbWFpbCI6Imh1YW5neXVhbjNoQGdtYWlsLmNvbSIsImV4cCI6MTcxOTM4MjU0NSwiaXNzIjoiaHR0cDovL25vcnRoLXBhdGguc2l0ZSIsInVzZXJOYW1lIjoi6buE57yY77yIWXVhbiBIdWFuZ--8iSJ9.YfAXmq3cq8pvx3A92WfHtfK2r9H1fUYGGL9Ts0eikTs",
		},
	}

	result, _ := Handler(input)

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", result.StatusCode)
	}
}
