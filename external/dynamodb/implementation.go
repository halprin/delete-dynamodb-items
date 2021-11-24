package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDb struct {
	service *dynamodb.DynamoDB
}

func NewDynamoDb() (*DynamoDb, error) {
	awsSession, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	dynamoService := dynamodb.New(awsSession)

	return &DynamoDb{
		service: dynamoService,
	}, nil
}

func NewDynamoDbWithEndpoint(endpoint string) (*DynamoDb, error) {
	awsSession, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	dynamoService := dynamodb.New(awsSession, aws.NewConfig().WithEndpoint(endpoint))

	return &DynamoDb{
		service: dynamoService,
	}, nil
}

func (d *DynamoDb) Describe(describeTableInput *dynamodb.DescribeTableInput) (*dynamodb.DescribeTableOutput, error) {
	panic("implement me")
}

func (d *DynamoDb) Scan(scanInput *dynamodb.ScanInput) chan []map[string]*dynamodb.AttributeValue {
	panic("implement me")
}

func (d *DynamoDb) BatchWrite(batchWriteItemInput *dynamodb.BatchWriteItemInput) *dynamodb.BatchWriteItemOutput {
	panic("implement me")
}

