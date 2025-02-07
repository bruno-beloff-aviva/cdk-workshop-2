package response

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type HelloResponse struct {
	StatusCode   int
	ErrorMessage string
	Body         string
}

func NewHelloResponse(statusCode int, err error, body string) HelloResponse {
	var errorMessage string

	if err == nil {
		errorMessage = ""
	} else {
		errorMessage = err.Error()
	}

	return HelloResponse{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
		Body:         body,
	}
}

func (r HelloResponse) APIResponse() (events.APIGatewayProxyResponse, error) {
	jsonBody, err := json.Marshal(r)

	return events.APIGatewayProxyResponse{
		StatusCode: r.StatusCode,
		Body:       string(jsonBody),
	}, err
}
