package response

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type HelloResponse struct {
	StatusCode   int
	ErrorMessage string
	Body         string
}

func NewOKHelloResponse(body string) HelloResponse {
	return HelloResponse{
		StatusCode:   http.StatusOK,
		ErrorMessage: "",
		Body:         body,
	}
}

func NewErrorHelloResponse(err error, body string) HelloResponse {
	return HelloResponse{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: err.Error(),
		Body:         body,
	}
}

func (r HelloResponse) APIResponse() (apiResponse events.APIGatewayProxyResponse, marshalErr error) {
	var jsonBody []byte

	jsonBody, marshalErr = json.Marshal(r)

	apiResponse = events.APIGatewayProxyResponse{
		StatusCode: r.StatusCode,
		Body:       string(jsonBody),
	}

	return apiResponse, marshalErr
}
