package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	fmt.Println("hitcounter_lambda init!!")

}

func handleRequest(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("hitcounter_lambda handleRequest!!")

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "hitcounter_lambda\n",
	}, nil
}

func main() {
	fmt.Println("hitcounter_lambda main!!")

	lambda.Start(handleRequest)
}
