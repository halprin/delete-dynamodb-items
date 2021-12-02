package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_getItemsGoroutine_filterExpressionIsNil(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())

	_ = getItemsGoroutine("a table name", nil, nil, nil)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.Nil(t, scanInputArgument.FilterExpression)
}

func Test_getItemsGoroutine_filterExpression(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())
	filterExpression := "a filter expression"

	_ = getItemsGoroutine("a table name", &filterExpression, nil, nil)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.Equal(t, filterExpression, *scanInputArgument.FilterExpression)
}

func Test_getItemsGoroutine_expressionAttributeNamesIsNil(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())

	_ = getItemsGoroutine("a table name", nil, nil, nil)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.Nil(t, scanInputArgument.ExpressionAttributeNames)
}

func Test_getItemsGoroutine_expressionAttributeNamesIsInvalid(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())
	expressionAttributeNames := "not valid JSON"

	goroutine := getItemsGoroutine("a table name", nil, &expressionAttributeNames, nil)

	assert.Nil(t, goroutine)
	mockDynamoDb.AssertNotCalled(t, "Scan", mock.Anything)
}

func Test_getItemsGoroutine_expressionAttributeNames(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())
	expressionAttributeNames := `{"dogcow": "moof"}`

	goroutine := getItemsGoroutine("a table name", nil, &expressionAttributeNames, nil)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.NotNil(t, goroutine)
	assert.NotNil(t, scanInputArgument.ExpressionAttributeNames)
}

func Test_getItemsGoroutine_expressionAttributeValuesIsNil(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())

	_ = getItemsGoroutine("a table name", nil, nil, nil)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.Nil(t, scanInputArgument.ExpressionAttributeValues)
}

func Test_getItemsGoroutine_expressionAttributeValuesIsInvalid(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())
	expressionAttributeValues := "not valid JSON"

	goroutine := getItemsGoroutine("a table name", nil, nil, &expressionAttributeValues)

	assert.Nil(t, goroutine)
	mockDynamoDb.AssertNotCalled(t, "Scan", mock.Anything)
}

func Test_getItemsGoroutine_expressionAttributeValues(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())
	expressionAttributeValues := `{"dogcow": {"BOOL": true}}`

	goroutine := getItemsGoroutine("a table name", nil, nil, &expressionAttributeValues)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.NotNil(t, goroutine)
	assert.NotNil(t, scanInputArgument.ExpressionAttributeValues)
}

func testChannelOfScanMethodReturnType() chan []map[string]*dynamodb.AttributeValue {
	channel := make(chan []map[string]*dynamodb.AttributeValue)

	go func() {
		arrayOfDynamoDbStuff := []map[string]*dynamodb.AttributeValue{
			{
				"dogcow": {
					S: aws.String("moof"),
				},
			},
		}
		channel <- arrayOfDynamoDbStuff
		close(channel)
	}()

	return channel
}
