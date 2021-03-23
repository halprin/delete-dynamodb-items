package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/halprin/delete-dynamodb-items/parallel"
	"log"
	"math/rand"
	"time"
)

var maxItemsPerBatchRequest = 25
var tableKeys []*dynamodb.KeySchemaElement

func deleteItems(dynamoItems []map[string]*dynamodb.AttributeValue, tableName string, goroutinePool *parallel.Pool) error {

	var err error
	tableKeys, err = getTableKeys(tableName)
	if err != nil {
		log.Println("Unable to determine the keys of the table")
		return err
	}

	dynamoItemsChunks := chunkItems(dynamoItems)

	var errorChannels []chan error

	for _, currentItemsChunk := range dynamoItemsChunks {

		errorChannel := make(chan error, 1)
		errorChannels = append(errorChannels, errorChannel)

		//wrapping in a function to make a copy of the currentItemsChunk and errorChannel arguments that are passed in,
		//else all executions try to delete the same chunk of items
		func(currentItemsChunk []map[string]*dynamodb.AttributeValue, errorChannel chan error) {
			goroutinePool.Submit(func() {
				deleteChunkGoroutine(currentItemsChunk, tableName, errorChannel)
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

	err := incrementallyBatchDelete(requestItems)
	if err != nil {
		log.Println("Failed to batch delete items")
		return err
	}

	return nil
}

func getTableKeys(tableName string) ([]*dynamodb.KeySchemaElement, error) {
	tableInfo, err := describeTable(tableName)
	if err != nil {
		return nil, err
	}
	return tableInfo.Table.KeySchema, nil
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

func incrementallyBatchDelete(requestItems map[string][]*dynamodb.WriteRequest) error {
	//used to induce jitter
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	baseMillisecondsToWait := 20
	maxMillisecondsToWait := 40
	millisecondsToWait := randomGenerator.Intn(maxMillisecondsToWait)

	//start of waiting so all the goroutines don't call batch delete at the same time
	time.Sleep(time.Duration(millisecondsToWait) * time.Millisecond)

	for {
		batchWriteItemInput := &dynamodb.BatchWriteItemInput{
			RequestItems: requestItems,
		}

		log.Println("Deleting some items")

		batchWriteItemOutput, err := dynamoService.BatchWriteItem(batchWriteItemInput)
		if err != nil {
			//there was an error writing to DynamoDB
			log.Println("Failed to put/delete items in DynamoDB")
			return err
		}

		if len(batchWriteItemOutput.UnprocessedItems) > 0 {
			//there are still items to write, reset requestItems for the next pass
			log.Println("Unprocessed items remain, trying again with remaining items")
			requestItems = batchWriteItemOutput.UnprocessedItems
		} else {
			//no more items to write, break out
			break
		}

		//do an exponential back-off with jitter
		time.Sleep(time.Duration(millisecondsToWait) * time.Millisecond)
		maxMillisecondsToWait *= 2
		millisecondsToWait = baseMillisecondsToWait + randomGenerator.Intn(maxMillisecondsToWait)
	}

	return nil
}
