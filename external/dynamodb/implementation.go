package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDb struct {
	service *dynamodb.DynamoDB
}

var singleton *DynamoDb

func InitializeDynamoDb() error {
	awsSession, err := session.NewSession()
	if err != nil {
		return err
	}

	dynamoService := dynamodb.New(awsSession)

	singleton = &DynamoDb{
		service: dynamoService,
	}

	return nil
}

func InitializeDynamoDbWithEndpoint(endpoint string) error {
	awsSession, err := session.NewSession()
	if err != nil {
		return err
	}

	dynamoService := dynamodb.New(awsSession, aws.NewConfig().WithEndpoint(endpoint))

	singleton = &DynamoDb{
		service: dynamoService,
	}

	return nil
}

func GetService() *DynamoDb {
	return singleton
}

func (d *DynamoDb) Describe(tableName string) (*dynamodb.DescribeTableOutput, error) {
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	tableInfo, err := d.service.DescribeTable(describeTableInput)
	if err != nil {
		return nil, err
	}

	return tableInfo, nil
}

func (d *DynamoDb) Scan(scanInput *dynamodb.ScanInput) chan []map[string]*dynamodb.AttributeValue {
	panic("implement me")
}

func (d *DynamoDb) BatchWrite(batchWriteItemInput *dynamodb.BatchWriteItemInput) *dynamodb.BatchWriteItemOutput {
	panic("implement me")
}

