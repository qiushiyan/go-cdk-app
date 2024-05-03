package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"golang.org/x/crypto/bcrypt"
)

const USERS_TABLE_NAME = "users"

type User struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
}

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *DynamoDBClient) UserExists(username string) (bool, error) {
	return c.exists(USERS_TABLE_NAME, username)
}

func (c *DynamoDBClient) InsertUser(nu NewUser) error {
	h, err := hashPassword(nu.Password)
	if err != nil {
		return err
	}

	_, err = c.insert(USERS_TABLE_NAME, map[string]*dynamodb.AttributeValue{
		"username": {
			S: aws.String(nu.Username),
		},
		"hashed_password": {
			S: aws.String(h),
		},
	})

	return err
}

func (c *DynamoDBClient) GetUser(username string) (User, error) {
	var user User

	output, err := c.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(USERS_TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return User{}, err
	}

	if output.Item == nil {
		return User{}, nil
	}

	err = dynamodbattribute.UnmarshalMap(output.Item, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func ValidatePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
