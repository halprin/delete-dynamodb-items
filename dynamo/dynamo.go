package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
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

func describeTable(tableName string) (*dynamodb.DescribeTableOutput, error) {
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	tableInfo, err := dynamoService.DescribeTable(describeTableInput)
	if err != nil {
		log.Println("Unable to describe the the table")
		return nil, err
	}

	return tableInfo, nil
}
