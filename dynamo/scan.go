package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

func getItems(tableName string) ([]map[string]*dynamodb.AttributeValue, error) {

	var scannedItems []map[string]*dynamodb.AttributeValue

	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	for {
		log.Println("Scanning items")

		scanOutput, err := dynamoService.Scan(scanInput)
		if err != nil {
			log.Println("Failed to scan the items")
			return nil, err
		}

		scannedItems = append(scannedItems, scanOutput.Items...)

		if scanOutput.LastEvaluatedKey != nil && len(scanOutput.LastEvaluatedKey) > 0 {
			//there are still items to scan, the the key to start scanning from again
			scanInput.ExclusiveStartKey = scanOutput.LastEvaluatedKey
		} else {
			//no more items to scan, break out
			break
		}
	}

	return scannedItems, nil
}

