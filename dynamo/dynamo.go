package dynamo

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

var awsSession, sessionErr = session.NewSession()
var dynamoService = dynamodb.New(awsSession)

func DeleteAllItemsInTable(tableName string) error {
	if sessionErr != nil {
		log.Println("Initial AWS session failed")
		return sessionErr
	}

	items, err := getItems(tableName)
	if err != nil {
		return err
	}

	err = deleteItems(items, tableName)
	return err
}
