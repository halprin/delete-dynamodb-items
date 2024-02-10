package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/halprin/delete-dynamodb-items/parallel"
	"log"
)

var maxItemsPerBatchRequest = 25
var tableKeys []types.KeySchemaElement

func deleteItems(dynamoItems []map[string]types.AttributeValue, tableInfo *dynamodb.DescribeTableOutput, goroutinePool *parallel.Pool) error {

	tableKeys = getTableKeys(tableInfo)

	dynamoItemsChunks := chunkItems(dynamoItems)

	var errorChannels []chan error

	for _, currentItemsChunk := range dynamoItemsChunks {

		errorChannel := make(chan error, 1)
		errorChannels = append(errorChannels, errorChannel)

		//wrapping in a function to make a copy of the currentItemsChunk and errorChannel arguments that are passed in,
		//else all executions try to delete the same chunk of items
		func(currentItemsChunk []map[string]types.AttributeValue, errorChannel chan error) {
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

func deleteChunkGoroutine(currentItemsChunk []map[string]types.AttributeValue, tableName string, errorChannel chan error) {
	errorChannel <- deleteChunk(currentItemsChunk, tableName)
	close(errorChannel)
}

func deleteChunk(currentItemsChunk []map[string]types.AttributeValue, tableName string) error {
	writeRequests := marshalItemsIntoBatchWrites(currentItemsChunk)

	requestItems := map[string][]types.WriteRequest{
		tableName: writeRequests,
	}

	err := GetService().BatchWrite(requestItems)
	if err != nil {
		log.Println("Failed to batch delete items")
		return err
	}

	return nil
}

func getTableKeys(tableInfo *dynamodb.DescribeTableOutput) []types.KeySchemaElement {
	return tableInfo.Table.KeySchema
}

func chunkItems(dynamoItems []map[string]types.AttributeValue) [][]map[string]types.AttributeValue {
	var itemChunks [][]map[string]types.AttributeValue
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

func marshalItemsIntoBatchWrites(dynamoItems []map[string]types.AttributeValue) []types.WriteRequest {
	var writeRequests []types.WriteRequest
	var writeRequest types.WriteRequest

	for _, currentDynamoItem := range dynamoItems {
		key := convertItemToKey(currentDynamoItem)

		deleteRequest := &types.DeleteRequest{
			Key: key,
		}

		writeRequest = types.WriteRequest{
			DeleteRequest: deleteRequest,
		}

		writeRequests = append(writeRequests, writeRequest)
	}

	return writeRequests
}

func convertItemToKey(item map[string]types.AttributeValue) map[string]types.AttributeValue {
	key := make(map[string]types.AttributeValue)
	for _, currentTableKey := range tableKeys {
		currentTableKeyName := *currentTableKey.AttributeName
		key[currentTableKeyName] = item[currentTableKeyName]
	}

	return key
}
