package dynamo

import (
	"encoding/base64"
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

	value, ok = rawAttributeValue["BOOL"]
	if ok {
		valueBool, ok := value.(bool)
		if !ok {
			return nil, errors.New("BOOL attribute value is not a boolean")
		}

		return &types.AttributeValueMemberBOOL{Value: valueBool}, nil
	}

	value, ok = rawAttributeValue["B"]
	if ok {
		valueString, ok := value.(string)
		if !ok {
			return nil, errors.New("B attribute value is not a string")
		}

		valueBytes, err := base64.StdEncoding.DecodeString(valueString)
		if err != nil {
			return nil, err
		}

		return &types.AttributeValueMemberB{Value: valueBytes}, nil
	}

	value, ok = rawAttributeValue["NULL"]
	if ok {
		return &types.AttributeValueMemberNULL{}, nil
	}

	value, ok = rawAttributeValue["SS"]
	if ok {
		valueList, ok := value.([]interface{})
		if !ok {
			return nil, errors.New("SS attribute value is not a list")
		}

		valueListOfStrings := make([]string, len(valueList))
		for index, item := range valueList {
			valueString, ok := item.(string)
			if !ok {
				return nil, errors.New("SS attribute value's sub-value is not a string")
			}
			valueListOfStrings[index] = valueString
		}

		return &types.AttributeValueMemberSS{Value: valueListOfStrings}, nil
	}

	value, ok = rawAttributeValue["NS"]
	if ok {
		valueList, ok := value.([]interface{})
		if !ok {
			return nil, errors.New("NS attribute value is not a list")
		}

		valueListOfStrings := make([]string, len(valueList))
		for index, item := range valueList {
			valueString, ok := item.(string)
			if !ok {
				return nil, errors.New("NS attribute value's sub-value is not a string")
			}
			valueListOfStrings[index] = valueString
		}

		return &types.AttributeValueMemberNS{Value: valueListOfStrings}, nil
	}

	value, ok = rawAttributeValue["BS"]
	if ok {
		valueListOfStrings, ok := value.([]string)
		if !ok {
			return nil, errors.New("BS attribute value is not a valueListOfBytes of strings")
		}

		valueListOfBytes := make([][]byte, len(valueListOfStrings))
		for index, item := range valueListOfStrings {
			bytes, err := base64.StdEncoding.DecodeString(item)
			if err != nil {
				return nil, err
			}

			valueListOfBytes[index] = bytes
		}

		return &types.AttributeValueMemberBS{Value: valueListOfBytes}, nil
	}

	value, ok = rawAttributeValue["L"]
	if ok {
		valueList, ok := value.([]map[string]interface{})
		if !ok {
			return nil, errors.New("L attribute value is not a list of maps from string to something")
		}

		list := make([]types.AttributeValue, len(valueList))
		for index, item := range valueList {
			attributeValue, err := convertRawAttributeValues(item)
			if err != nil {
				return nil, err
			}

			list[index] = attributeValue
		}

		return &types.AttributeValueMemberL{Value: list}, nil
	}

	value, ok = rawAttributeValue["M"]
	if ok {
		valueMap, ok := value.(map[string]map[string]interface{})
		if !ok {
			return nil, errors.New("M attribute value is not a map")
		}

		mappedAttributeValues, err := mapAttributeValues(valueMap)
		if err != nil {
			return nil, err
		}

		return &types.AttributeValueMemberM{Value: mappedAttributeValues}, nil
	}

	return nil, errors.New("attribute value type didn't match any known types")
}
