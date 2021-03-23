package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/halprin/delete-dynamodb-items/config"
	"github.com/halprin/delete-dynamodb-items/parallel"
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

	concurrency, err := determineConcurrency(tableName)
	if err != nil {
		log.Println("Unable determine the concurrency")
		return err
	}

	// 1024 * 1024 / 25 = 41,943.04 ~= 41,944
	goroutinePool := parallel.NewPool(concurrency, 41944)
	defer goroutinePool.Release()

	for subItemList := range getItemsGoroutine(tableName) {
		err = deleteItems(subItemList, tableName, goroutinePool)
		if err != nil {
			return err
		}
	}

	return nil
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
