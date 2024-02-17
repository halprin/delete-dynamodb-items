package dynamo

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"log"
)

type expressionAttributeNamesType map[string]string
type expressionAttributeValuesType map[string]types.AttributeValue

//type expressionAttributeValuesType types.AttributeValueMemberM

func getItemsGoroutine(tableName string, filterExpression *string, expressionAttributeNames *string, expressionAttributeValues *string) chan []map[string]types.AttributeValue {

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
		var mappedData map[string]map[string]interface{}
		err := json.Unmarshal([]byte(*expressionAttributeValues), &mappedData)
		if err != nil {
			log.Printf("Failed to unmarshal the expression attribute values, %+v", err)
			return nil
		}

		values, err = mapAttributeValues(mappedData)
		if err != nil {
			log.Printf("Failed to map the expression attribute values, %+v", err)
			return nil
		}
	}

	scanInput := &dynamodb.ScanInput{
		TableName:                 &tableName,
		FilterExpression:          filterExpression,
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
	}

	return GetService().Scan(scanInput)
}

func mapAttributeValues(unmarshaledJson map[string]map[string]interface{}) (expressionAttributeValuesType, error) {
	expressionAttributeValues := expressionAttributeValuesType{}
	for key, rawAttributeValue := range unmarshaledJson {
		attributeValue, err := convertRawAttributeValues(rawAttributeValue)
		if err != nil {
			return nil, err
		}
		expressionAttributeValues[key] = attributeValue
	}

	return expressionAttributeValues, nil
}

func convertRawAttributeValues(rawAttributeValue map[string]interface{}) (types.AttributeValue, error) {

	value, ok := rawAttributeValue["N"]
	if ok {
		valueString, ok := value.(string)
		if !ok {
			return nil, errors.New("N attribute value is not a string")
		}
		return &types.AttributeValueMemberN{Value: valueString}, nil
	}

	value, ok = rawAttributeValue["S"]
	if ok {
		valueString, ok := value.(string)
		if !ok {
			return nil, errors.New("S attribute value is not a string")
		}
		return &types.AttributeValueMemberS{Value: valueString}, nil
	}

	return nil, errors.New("attribute value type didn't match any known types")
}
