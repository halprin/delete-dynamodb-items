package dynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
	"math/rand"
	"time"
)

type DynamoDb struct {
	client *dynamodb.Client
}

func NewDynamoDb() (*DynamoDb, error) {
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	dynamoClient := dynamodb.NewFromConfig(awsConfig)

	return &DynamoDb{
		client: dynamoClient,
	}, nil
}

func NewDynamoDbWithEndpoint(endpoint string) (*DynamoDb, error) {
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	dynamoClient := dynamodb.NewFromConfig(awsConfig, func(options *dynamodb.Options) {
		options.BaseEndpoint = &endpoint
	})

	return &DynamoDb{
		client: dynamoClient,
	}, nil
}

func (d *DynamoDb) Describe(tableName string) (*dynamodb.DescribeTableOutput, error) {
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	tableInfo, err := d.client.DescribeTable(context.Background(), describeTableInput)
	if err != nil {
		return nil, err
	}

	return tableInfo, nil
}

func (d *DynamoDb) Scan(scanInput *dynamodb.ScanInput) chan []map[string]types.AttributeValue {
	yield := make(chan []map[string]types.AttributeValue)

	go func() {

		for {
			log.Println("Scanning items")

			scanOutput, err := d.client.Scan(context.Background(), scanInput)
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

func (d *DynamoDb) BatchWrite(requestItems map[string][]types.WriteRequest) error {
	//used to induce jitter
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	baseMillisecondsToWait := 20
	maxMillisecondsToWait := 40
	millisecondsToWait := randomGenerator.Intn(maxMillisecondsToWait)

	//start of waiting so all the goroutines don't call batch delete at the same time
	time.Sleep(time.Duration(millisecondsToWait) * time.Millisecond)

	for {
		batchWriteItemInput := &dynamodb.BatchWriteItemInput{
			RequestItems: requestItems,
		}

		log.Println("Deleting some items")

		batchWriteItemOutput, err := d.client.BatchWriteItem(context.Background(), batchWriteItemInput)
		if err != nil {
			//there was an error writing to DynamoDB
			log.Println("Failed to put/delete items in DynamoDB")
			return err
		}

		if len(batchWriteItemOutput.UnprocessedItems) > 0 {
			//there are still items to write, reset requestItems for the next pass
			log.Println("Unprocessed items remain, trying again with remaining items")
			requestItems = batchWriteItemOutput.UnprocessedItems
		} else {
			//no more items to write, break out
			break
		}

		//do an exponential back-off with jitter
		time.Sleep(time.Duration(millisecondsToWait) * time.Millisecond)
		maxMillisecondsToWait *= 2
		millisecondsToWait = baseMillisecondsToWait + randomGenerator.Intn(maxMillisecondsToWait)
	}

	return nil
}
