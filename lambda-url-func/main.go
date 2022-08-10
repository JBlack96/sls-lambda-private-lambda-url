package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	fmt.Println("Hello from Private Lambda URL Function!")
	fmt.Println("Looking to fix up a response")
	return LambdaResponse(http.StatusOK, "looks mint like"), nil
}

func main() {
	lambda.Start(Handler)
}

// LambdaResponse is a helper function to create a LambdaFunctionURLResponse from a status code and a body.
func LambdaResponse(status int, body string) events.LambdaFunctionURLResponse {
	resp := events.LambdaFunctionURLResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}

	resp.Body = body
	return resp
}
