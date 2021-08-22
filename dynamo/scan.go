package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

func getItemsGoroutine(tableName string) chan []map[string]*dynamodb.AttributeValue {
	yield := make(chan []map[string]*dynamodb.AttributeValue)

	go func() {
		scanInput := &dynamodb.ScanInput{
			TableName:                aws.String(tableName),
			FilterExpression:         aws.String("#k > :v"),
			ExpressionAttributeNames: map[string]*string{
				"#k": aws.String("number"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":v": {
					N: aws.String("50"),
				},
			},
		}

		for {
			log.Println("Scanning items")

			scanOutput, err := dynamoService.Scan(scanInput)
			if err != nil {
				log.Printf("Failed to scan the items, %+v", err)
				break
			}

			yield <- scanOutput.Items

			if scanOutput.LastEvaluatedKey != nil && len(scanOutput.LastEvaluatedKey) > 0 {
				//there are still items to scan, the the key to start scanning from again
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
