package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/halprin/delete-dynamodb-items/parallel"
	"log"
)

var maxItemsPerBatchRequest = 25
var tableKeys []*dynamodb.KeySchemaElement

func deleteItems(dynamoItems []map[string]*dynamodb.AttributeValue, tableInfo *dynamodb.DescribeTableOutput, goroutinePool *parallel.Pool) error {

	tableKeys = getTableKeys(tableInfo)

	dynamoItemsChunks := chunkItems(dynamoItems)

	var errorChannels []chan error

	for _, currentItemsChunk := range dynamoItemsChunks {

		errorChannel := make(chan error, 1)
		errorChannels = append(errorChannels, errorChannel)

		//wrapping in a function to make a copy of the currentItemsChunk and errorChannel arguments that are passed in,
		//else all executions try to delete the same chunk of items
		func(currentItemsChunk []map[string]*dynamodb.AttributeValue, errorChannel chan error) {
			goroutinePool.Submit(func() {
				deleteChunkGoroutine(currentItemsChunk, *tableInfo.Table.TableName, errorChannel)
			})
		}(currentItemsChunk, errorChannel)
	}

	log.Println("Waiting for all deletion goroutines to complete")

	for errorFromGoroutine := range parallel.MergeErrorChannels(errorChannels) {
		if errorFromGoroutine != nil {
			log.Println("One of the delete goroutines failed")
			return errorFromGoroutine
		}
	}

	return nil
}

func deleteChunkGoroutine(currentItemsChunk []map[string]*dynamodb.AttributeValue, tableName string, errorChannel chan error) {
	errorChannel <- deleteChunk(currentItemsChunk, tableName)
	close(errorChannel)
}

func deleteChunk(currentItemsChunk []map[string]*dynamodb.AttributeValue, tableName string) error {
	writeRequests := marshalItemsIntoBatchWrites(currentItemsChunk)

	requestItems := map[string][]*dynamodb.WriteRequest{
		tableName: writeRequests,
	}

	err := GetService().BatchWrite(requestItems)
	if err != nil {
		log.Println("Failed to batch delete items")
		return err
	}

	return nil
}

func getTableKeys(tableInfo *dynamodb.DescribeTableOutput) []*dynamodb.KeySchemaElement {
	return tableInfo.Table.KeySchema
}

func chunkItems(dynamoItems []map[string]*dynamodb.AttributeValue) [][]map[string]*dynamodb.AttributeValue {
	var itemChunks [][]map[string]*dynamodb.AttributeValue
	numberOfItems := len(dynamoItems)

	for itemIndex := 0; itemIndex < numberOfItems; itemIndex += maxItemsPerBatchRequest {
		end := itemIndex + maxItemsPerBatchRequest

		if end > numberOfItems {
			end = numberOfItems
		}

		itemChunks = append(itemChunks, dynamoItems[itemIndex:end])
	}

	return itemChunks
}

func marshalItemsIntoBatchWrites(dynamoItems []map[string]*dynamodb.AttributeValue) []*dynamodb.WriteRequest {
	var writeRequests []*dynamodb.WriteRequest
	var writeRequest *dynamodb.WriteRequest

	for _, currentDynamoItem := range dynamoItems {
		key := convertItemToKey(currentDynamoItem)

		deleteRequest := &dynamodb.DeleteRequest{
			Key:  key,
		}

		writeRequest = &dynamodb.WriteRequest{
			DeleteRequest: deleteRequest,
		}

		writeRequests = append(writeRequests, writeRequest)
	}

	return writeRequests
}

func convertItemToKey(item map[string]*dynamodb.AttributeValue) map[string]*dynamodb.AttributeValue {
	key := make(map[string]*dynamodb.AttributeValue)
	for _, currentTableKey := range tableKeys {
		currentTableKeyName := *currentTableKey.AttributeName
		key[currentTableKeyName] = item[currentTableKeyName]
	}

	return key
}
