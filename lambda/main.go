package main

import (
	"go-aws/lambda/app"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	app := app.New()
	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch req.Path {
		case "/register":
			return app.Handler.RegisterUser(req)
		case "/login":
			return app.Handler.LoginUser(req)
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       "Not Found",
			}, nil
		}
	})
}
