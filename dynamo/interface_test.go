package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/mock"
)

func ResetDynamoDbMock() *DynamoDbMock {
	dynamoDbMock := &DynamoDbMock{}
	singleton = dynamoDbMock
	return dynamoDbMock
}

type DynamoDbMock struct {
	mock.Mock
}

func (d *DynamoDbMock) Describe(tableName string) (*dynamodb.DescribeTableOutput, error) {
	args := d.Called(tableName)
	return args.Get(0).(*dynamodb.DescribeTableOutput), args.Error(1)
}

func (d *DynamoDbMock) Scan(input *dynamodb.ScanInput) chan []map[string]*dynamodb.AttributeValue {
	args := d.Called(input)
	return args.Get(0).(chan []map[string]*dynamodb.AttributeValue)
}

func (d *DynamoDbMock) BatchWrite(requestItems map[string][]*dynamodb.WriteRequest) error {
	args := d.Called(requestItems)
	return args.Error(0)
}
