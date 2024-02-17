package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	myDynamoImplementation "github.com/halprin/delete-dynamodb-items/external/dynamodb"
)

type DynamoDber interface {
	Describe(tableName string) (*dynamodb.DescribeTableOutput, error)
	Scan(*dynamodb.ScanInput) chan []map[string]types.AttributeValue
	BatchWrite(requestItems map[string][]types.WriteRequest) error
}

var singleton DynamoDber

func InitializeDynamoDb() error {
	if singleton != nil {
		panic("singleton already initialized, call GetService()")
	}

	service, err := myDynamoImplementation.NewDynamoDb()
	if err != nil {
		return err
	}

	singleton = service
	return nil
}

func InitializeDynamoDbWithEndpoint(endpoint string) error {
	if singleton != nil {
		panic("singleton already initialized, call GetService()")
	}

	service, err := myDynamoImplementation.NewDynamoDbWithEndpoint(endpoint)
	if err != nil {
		return err
	}

	singleton = service
	return nil
}

func GetService() DynamoDber {
	return singleton
}
