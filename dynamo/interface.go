package dynamo

import "github.com/aws/aws-sdk-go/service/dynamodb"

type DynamoDber interface {
	Describe(tableName string) (*dynamodb.DescribeTableOutput, error)
	Scan(*dynamodb.ScanInput) chan []map[string]*dynamodb.AttributeValue
	BatchWrite(requestItems map[string][]*dynamodb.WriteRequest) error
}
