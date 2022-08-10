package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler provides deps needed for the handle func
type Handler struct {
	signer *v4.Signer
}

// Handle is our lambda handler invoked by the `lambda.Start` function call
func (h *Handler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	fmt.Println("getting lambda url function URL")
	// check env for url
	url := os.Getenv("PRIVATE_LAMBDA_URL")
	if url == "" {
		fmt.Println("failed to get url from env")
		return writeResponse("Failed while executing", 500)
	}
	fmt.Println("successfully got URL:", url)

	fmt.Println("creating http client")
	client := http.Client{
		Timeout: time.Second * 5,
	}
	fmt.Println("created client")

	fmt.Println("creating request")
	req, err := http.NewRequest("POST", url, nil)
	fmt.Println("created request")

	if err != nil {
		fmt.Println("failed to call format request")
		return writeResponse("Failed while executing", 500)
	}

	fmt.Println("sign the request with aws v4")
	_, err = h.signer.Sign(req, nil, "lambda", "us-east-1", time.Now())
	if err != nil {
		fmt.Println("failed to generate aws sigv4 for http request")
		return writeResponse("Failed while executing", 500)
	}
	fmt.Println("signed the request with aws v4")

	fmt.Println("calling lambda url function")
	resp, err := client.Do(req)
	fmt.Println("called lambda url function")

	if err != nil {
		fmt.Println("failed to call lambda url successfully")
		return writeResponse("Failed while executing", 500)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("didn't receive a 200 but got", resp.StatusCode)
		return writeResponse("Failed while executing", 500)
	}

	return writeResponse("Successfully executed", 200)
}

func writeResponse(msg string, status int) (Response, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"message": msg,
	})
	resp := Response{
		StatusCode:      status,
		IsBase64Encoded: false,
		Body:            string(body),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "entry-handler",
		},
	}
	return resp, nil
}

func main() {
	sess := session.Must(session.NewSession())
	signer := v4.NewSigner(sess.Config.Credentials)

	handler := Handler{
		signer: signer,
	}

	lambda.Start(handler.Handle)
}
