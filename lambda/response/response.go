package response

import "github.com/aws/aws-lambda-go/events"

func New200(message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       message + "\n",
	}
}
