package dynamo

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/halprin/delete-dynamodb-items/parallel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_chunkItems(t *testing.T) {
	chunkedDynamoDbItems := chunkItems(testDynamoDbItems())

	for _, dynamoDbItems := range chunkedDynamoDbItems {
		assert.LessOrEqual(t, len(dynamoDbItems), maxItemsPerBatchRequest)
	}
}

func Test_deleteChunk_succeeds(t *testing.T) {
	assert := assert.New(t)

	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("BatchWrite", mock.Anything).Return(nil)
	tableName := "DogCow"
	columnName := "moofColumn"
	dynamoDbItems := testDynamoDbItems(columnName)
	tableKeys = []*dynamodb.KeySchemaElement{
		{
			AttributeName: &columnName,
		},
	}

	err := deleteChunk(dynamoDbItems, tableName)

	writeRequestArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(map[string][]*dynamodb.WriteRequest)

	assert.Nil(err)
	assert.Len(writeRequestArgument[tableName], len(dynamoDbItems))
	for _, request := range writeRequestArgument[tableName] {
		assert.NotNil(request.DeleteRequest)
		assert.NotEmpty(request.DeleteRequest.Key)
	}
}

func Test_deleteChunk_errorsFromCallToDynamoDb(t *testing.T) {
	assert := assert.New(t)

	mockDynamoDb := ResetDynamoDbMock()
	mockError := errors.New("this is an error")
	mockDynamoDb.On("BatchWrite", mock.Anything).Return(mockError)
	tableName := "DogCow"
	columnName := "moofColumn"
	dynamoDbItems := testDynamoDbItems(columnName)
	tableKeys = []*dynamodb.KeySchemaElement{
		{
			AttributeName: &columnName,
		},
	}

	err := deleteChunk(dynamoDbItems, tableName)

	assert.Equal(mockError, err)
}

func Test_deleteItems_succeeds(t *testing.T) {
	assert := assert.New(t)

	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("BatchWrite", mock.Anything).Return(nil)
	tableName := "DogCow"
	columnName := "moofColumn"
	tableInfo := &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			TableName: &tableName,
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: &columnName,
				},
			},
		},
	}
	dynamoDbItems := testDynamoDbItems(columnName)

	//create a simple goroutine pool
	goroutinePool := parallel.NewPool(1, 50)
	defer goroutinePool.Release()

	err := deleteItems(dynamoDbItems, tableInfo, goroutinePool)

	assert.Nil(err)
	for _, mockCall := range mockDynamoDb.Calls {
		writeRequestArgument := mockCall.Arguments.Get(0).(map[string][]*dynamodb.WriteRequest)
		for _, request := range writeRequestArgument[tableName] {
			keyToBeDeleted := request.DeleteRequest.Key
			//remove item from the dynamoDB items.
			index := indexOf(dynamoDbItems, keyToBeDeleted)
			dynamoDbItems = append(dynamoDbItems[:index], dynamoDbItems[index+1:]...)
		}
	}
	assert.Empty(dynamoDbItems) //dynamoDbItems being empty means we deleted every item passed into deleteItems
}

func indexOf(slice []map[string]*dynamodb.AttributeValue, valueToSearchFor map[string]*dynamodb.AttributeValue) int {
	for index, valueInSlice := range slice {
		if mapsEqual(valueInSlice, valueToSearchFor) {
			return index
		}
	}

	return -1
}

func mapsEqual(map1 map[string]*dynamodb.AttributeValue, map2 map[string]*dynamodb.AttributeValue) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key1, value1 := range map1 {
		if map2[key1] != value1 {
			return false
		}
	}

	return true
}

func testDynamoDbItems(columnName ...string) []map[string]*dynamodb.AttributeValue {
	sliceCapacity := 128
	columnNameToUse := "dogcowColumn"
	if len(columnName) > 0 {
		columnNameToUse = columnName[0]
	}

	sliceOfDynamoDbitems := make([]map[string]*dynamodb.AttributeValue, 0, sliceCapacity)

	for itemIndex := 0; itemIndex < sliceCapacity; itemIndex++ {
		sliceOfDynamoDbitems = append(sliceOfDynamoDbitems, map[string]*dynamodb.AttributeValue{
			columnNameToUse: {
				S: aws.String(fmt.Sprintf("moof%d", itemIndex)),
			},
		})
	}

	return sliceOfDynamoDbitems
}
