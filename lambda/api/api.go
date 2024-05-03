package api

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

func decode(req events.APIGatewayProxyRequest, v interface{}) error {
	if err := json.Unmarshal([]byte(req.Body), v); err != nil {
		return fmt.Errorf("could not unmarshal request: %v", err)
	}
	return nil
}
