package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)


func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Received GitHub webhook event")

	// We'll implement the actual logic later

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Processing PR review request",
	}, nil
}

func main() {
	lambda.Start(Handler)
}
