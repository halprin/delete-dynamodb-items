package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
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
	yield := make(chan []map[string]*dynamodb.AttributeValue)

	go func() {

		for {
			log.Println("Scanning items")

			scanOutput, err := d.service.Scan(scanInput)
			if err != nil {
				log.Printf("Failed to scan the items, %+v", err)
				break
			}

			yield <- scanOutput.Items

			if scanOutput.LastEvaluatedKey != nil && len(scanOutput.LastEvaluatedKey) > 0 {
				//there are still items to scan, set the key to start scanning from again
				scanInput.ExclusiveStartKey = scanOutput.LastEvaluatedKey
			} else {
				//no more items to scan, break out
				break
			}
		}
		close(yield)
	}()

	return yield
}

func (d *DynamoDb) BatchWrite(batchWriteItemInput *dynamodb.BatchWriteItemInput) *dynamodb.BatchWriteItemOutput {
	panic("implement me")
}

