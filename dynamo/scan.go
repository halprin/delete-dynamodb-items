package dynamo

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
)

type expressionAttributeNamesType map[string]*string
type expressionAttributeValuesType map[string]*dynamodb.AttributeValue

func getItemsGoroutine(tableName string, filterExpression *string, expressionAttributeNames *string, expressionAttributeValues *string) chan []map[string]*dynamodb.AttributeValue {

	var names expressionAttributeNamesType
	if expressionAttributeNames != nil {
		err := json.Unmarshal([]byte(*expressionAttributeNames), &names)
		if err != nil {
			log.Printf("Failed to unmarshal the expression attribute names, %+v", err)
			return nil
		}
	}

	var values expressionAttributeValuesType
	if expressionAttributeValues != nil {
		err := json.Unmarshal([]byte(*expressionAttributeValues), &values)
		if err != nil {
			log.Printf("Failed to unmarshal the expression attribute values, %+v", err)
			return nil
		}
	}

	scanInput := &dynamodb.ScanInput{
		TableName:                aws.String(tableName),
		FilterExpression:         filterExpression,
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
	}

	return GetService().Scan(scanInput)
}
