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
