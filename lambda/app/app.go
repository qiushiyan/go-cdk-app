package app

import (
	"go-aws/lambda/api"
	"go-aws/lambda/database"

	"github.com/aws/aws-sdk-go/aws/session"
)

type App struct {
	Handler *api.Handler
}

func New() *App {
	s := session.Must(session.NewSession())

	db := database.NewDynamoDBClient(s)
	handler := api.NewHandler(db)
	return &App{Handler: handler}
}
