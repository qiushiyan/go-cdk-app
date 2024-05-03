package api

import "go-aws/lambda/database"

type Store interface {
	UserExists(username string) (bool, error)
	InsertUser(nu database.NewUser) error
	GetUser(string) (database.User, error)
}
