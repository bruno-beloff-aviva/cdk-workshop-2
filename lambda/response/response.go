package response

import "github.com/aws/aws-lambda-go/events"

func New(statusCode int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body + "\n",
	}
}
