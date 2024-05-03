package api

import (
	"errors"
	"fmt"
	"go-aws/lambda/database"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var newResponse = func(body string, status int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
	}
}

type EventRegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type EventLoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) RegisterUser(
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	var event EventRegisterUser
	if err := decode(req, &event); err != nil {
		return newResponse("invalid body", http.StatusBadRequest), err
	}

	username := event.Username
	password := event.Password

	if username == "" || password == "" {
		return newResponse("username or password is empty", http.StatusBadRequest), errors.New(
			"username and password are required",
		)
	}

	exists, err := h.store.UserExists(username)
	if err != nil {
		return newResponse("internal error", http.StatusInternalServerError), err
	}

	if exists {
		return events.APIGatewayProxyResponse{
				StatusCode: http.StatusConflict,
				Body:       "user already exists",
			}, fmt.Errorf(
				"user already exists with name %s",
				username,
			)
	}

	err = h.store.InsertUser(database.NewUser{
		Username: username,
		Password: password,
	})

	if err != nil {
		return newResponse("internal error", http.StatusInternalServerError), err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "user created",
	}, nil
}

func (h *Handler) LoginUser(
	req events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	var event EventLoginUser
	if err := decode(req, &event); err != nil {
		return newResponse(
				"invalid body",
				http.StatusBadRequest,
			), fmt.Errorf(
				"could not unmarshal request: %v",
				err,
			)
	}

	username := event.Username
	password := event.Password

	user, err := h.store.GetUser(username)
	if err != nil {
		return newResponse("internal error", http.StatusInternalServerError), err
	}

	if !database.ValidatePassword(password, user.HashedPassword) {
		return newResponse(
				"incorrect password",
				http.StatusUnauthorized,
			), errors.New(
				"incorrect password",
			)
	}

	return newResponse("login successful", http.StatusOK), nil
}
