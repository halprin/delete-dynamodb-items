package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/halprin/delete-dynamodb-items/config"
	"log"
)

var awsSession, sessionErr = session.NewSession()
var dynamoService = dynamodb.New(awsSession)

func DeleteAllItemsInTable() error {
	if sessionErr != nil {
		log.Println("Initial AWS session failed")
		return sessionErr
	}

	endpoint := config.GetDynamoDbEndpoint()
	if endpoint != nil {
		log.Printf("Using the custom endpoint %s", *endpoint)
		dynamoService = dynamodb.New(awsSession, aws.NewConfig().WithEndpoint(*endpoint))
	}

	tableName := *config.GetTableName()

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
