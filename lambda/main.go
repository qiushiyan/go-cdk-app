package main

import (
	"errors"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", errors.New("username is empty")
	}

	return "Hello " + event.Username, nil
}

func main() {
	lambda.Start(HandleRequest)
}
