package dynamo

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
