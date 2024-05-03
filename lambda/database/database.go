package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDBClient struct {
	db *dynamodb.DynamoDB
}

func NewDynamoDBClient(s *session.Session) *DynamoDBClient {
	db := dynamodb.New(s)
	return &DynamoDBClient{
		db: db,
	}
}

func (c *DynamoDBClient) exists(tableName, key string) (bool, error) {
	output, err := c.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		return false, err
	}

	if output.Item == nil {
		return false, nil
	}

	return true, nil
}

func (c *DynamoDBClient) insert(
	tableName string,
	item map[string]*dynamodb.AttributeValue,
) (*dynamodb.PutItemOutput, error) {
	return c.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
}
