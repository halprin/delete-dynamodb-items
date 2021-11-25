package dynamo

import (
	"github.com/halprin/delete-dynamodb-items/config"
	"github.com/halprin/delete-dynamodb-items/parallel"
	"log"
)

func DeleteAllItemsInTable() error {
	var err error

	endpoint := config.GetDynamoDbEndpoint()
	if endpoint == nil {
		err = InitializeDynamoDb()
	} else {
		log.Printf("Using the custom endpoint %s", *endpoint)
		err = InitializeDynamoDbWithEndpoint(*endpoint)
	}

	if err != nil {
		log.Println("Initial AWS session failed")
		return err
	}

	tableName := *config.GetTableName()

	tableInfo, err := GetService().Describe(tableName)
	if err != nil {
		log.Println("Unable to describe the the table")
		return err
	}

	concurrency := determineConcurrency(tableInfo)

	// 1024 * 1024 / 25 = 41,943.04 ~= 41,944
	goroutinePool := parallel.NewPool(concurrency, 41944)
	defer goroutinePool.Release()

	expressionFilter := config.GetFilterExpression()
	expressionAttributeNames := config.GetExpressionAttributeNames()
	expressionAttributeValues := config.GetExpressionAttributeValues()

	for subItemList := range getItemsGoroutine(tableName, expressionFilter, expressionAttributeNames, expressionAttributeValues) {
		err = deleteItems(subItemList, tableInfo, goroutinePool)
		if err != nil {
			return err
		}
	}

	return nil
}
