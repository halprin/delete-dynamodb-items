package dynamo

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_getItemsGoroutine_filterExpressionIsNil(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())

	_ = getItemsGoroutine("a table name", nil, nil, nil)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.Nil(t, scanInputArgument.FilterExpression)
}

func Test_getItemsGoroutine_filterExpression(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())
	filterExpression := "a filter expression"

	_ = getItemsGoroutine("a table name", &filterExpression, nil, nil)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.Equal(t, filterExpression, *scanInputArgument.FilterExpression)
}

func Test_getItemsGoroutine_expressionAttributeNamesIsNil(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())

	_ = getItemsGoroutine("a table name", nil, nil, nil)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.Nil(t, scanInputArgument.ExpressionAttributeNames)
}

func Test_getItemsGoroutine_expressionAttributeNamesIsInvalid(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())
	expressionAttributeNames := "not valid JSON"

	goroutine := getItemsGoroutine("a table name", nil, &expressionAttributeNames, nil)

	assert.Nil(t, goroutine)
	mockDynamoDb.AssertNotCalled(t, "Scan", mock.Anything)
}

func Test_getItemsGoroutine_expressionAttributeNames(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())
	expressionAttributeNames := `{"dogcow": "moof"}`

	goroutine := getItemsGoroutine("a table name", nil, &expressionAttributeNames, nil)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.NotNil(t, goroutine)
	assert.NotNil(t, scanInputArgument.ExpressionAttributeNames)
}

func Test_getItemsGoroutine_expressionAttributeValuesIsNil(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())

	_ = getItemsGoroutine("a table name", nil, nil, nil)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.Nil(t, scanInputArgument.ExpressionAttributeValues)
}

func Test_getItemsGoroutine_expressionAttributeValuesIsInvalid(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())
	expressionAttributeValues := "not valid JSON"

	goroutine := getItemsGoroutine("a table name", nil, nil, &expressionAttributeValues)

	assert.Nil(t, goroutine)
	mockDynamoDb.AssertNotCalled(t, "Scan", mock.Anything)
}

func Test_getItemsGoroutine_expressionAttributeValues(t *testing.T) {
	mockDynamoDb := ResetDynamoDbMock()
	mockDynamoDb.On("Scan", mock.Anything).Return(testChannelOfScanMethodReturnType())
	expressionAttributeValues := `{"dogcow": {"BOOL": true}}`

	goroutine := getItemsGoroutine("a table name", nil, nil, &expressionAttributeValues)

	scanInputArgument := mockDynamoDb.Calls[0].Arguments.Get(0).(*dynamodb.ScanInput)

	assert.NotNil(t, goroutine)
	assert.NotNil(t, scanInputArgument.ExpressionAttributeValues)
}

func Test_convertRawAttributeValues_N(t *testing.T) {
	expressionAttributeString := `{"N": "123.45"}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	attributeValue, err := convertRawAttributeValues(rawAttributeValue)

	assert.NoError(t, err)
	assert.NotNil(t, attributeValue)
}

func Test_convertRawAttributeValues_N_notString(t *testing.T) {
	expressionAttributeString := `{"N": 123.45}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_S(t *testing.T) {
	expressionAttributeString := `{"S": "Hello"}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	attributeValue, err := convertRawAttributeValues(rawAttributeValue)

	assert.NoError(t, err)
	assert.NotNil(t, attributeValue)
}

func Test_convertRawAttributeValues_S_notString(t *testing.T) {
	expressionAttributeString := `{"S": true}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_BOOL(t *testing.T) {
	expressionAttributeString := `{"BOOL": true}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	attributeValue, err := convertRawAttributeValues(rawAttributeValue)

	assert.NoError(t, err)
	assert.NotNil(t, attributeValue)
}

func Test_convertRawAttributeValues_BOOL_notBoolean(t *testing.T) {
	expressionAttributeString := `{"BOOL": "true"}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_B(t *testing.T) {
	expressionAttributeString := `{"B": "dGhpcyB0ZXh0IGlzIGJhc2U2NC1lbmNvZGVk"}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	attributeValue, err := convertRawAttributeValues(rawAttributeValue)

	assert.NoError(t, err)
	assert.NotNil(t, attributeValue)
}

func Test_convertRawAttributeValues_B_notString(t *testing.T) {
	expressionAttributeString := `{"B": 26.9}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_B_notBase64(t *testing.T) {
	expressionAttributeString := `{"B": "DogCow"}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_NULL(t *testing.T) {
	expressionAttributeString := `{"NULL": true}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	attributeValue, err := convertRawAttributeValues(rawAttributeValue)

	assert.NoError(t, err)
	assert.NotNil(t, attributeValue)
}

func Test_convertRawAttributeValues_SS(t *testing.T) {
	expressionAttributeString := `{"SS": ["Giraffe", "Hippo", "Zebra"]}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	attributeValue, err := convertRawAttributeValues(rawAttributeValue)

	assert.NoError(t, err)
	assert.NotNil(t, attributeValue)
}

func Test_convertRawAttributeValues_SS_notList(t *testing.T) {
	expressionAttributeString := `{"SS": "DogCow"}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_SS_notListOfStrings(t *testing.T) {
	expressionAttributeString := `{"SS": ["Giraffe", 26.9 ,"Zebra"]}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_NS(t *testing.T) {
	expressionAttributeString := `{"NS": ["42.2", "-19", "7.5", "3.14"]}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	attributeValue, err := convertRawAttributeValues(rawAttributeValue)

	assert.NoError(t, err)
	assert.NotNil(t, attributeValue)
}

func Test_convertRawAttributeValues_NS_notList(t *testing.T) {
	expressionAttributeString := `{"NS": "DogCow"}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_NS_notListOfStrings(t *testing.T) {
	expressionAttributeString := `{"NS": ["42.2", "-19", 7.5, "3.14"]}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_BS(t *testing.T) {
	expressionAttributeString := `{"BS": ["U3Vubnk=", "UmFpbnk=", "U25vd3k="]}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	attributeValue, err := convertRawAttributeValues(rawAttributeValue)

	assert.NoError(t, err)
	assert.NotNil(t, attributeValue)
}

func Test_convertRawAttributeValues_BS_notList(t *testing.T) {
	expressionAttributeString := `{"BS": "DogCow"}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_BS_notListOfStrings(t *testing.T) {
	expressionAttributeString := `{"BS": ["U3Vubnk=", 26.9, "U25vd3k="]}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_BS_notBase64(t *testing.T) {
	expressionAttributeString := `{"BS": ["U3Vubnk=", "DogCow", "U25vd3k="]}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_L(t *testing.T) {
	expressionAttributeString := `{"L": [ {"S": "Cookies"} , {"S": "Coffee"}, {"N": "3.14159"}]}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	attributeValue, err := convertRawAttributeValues(rawAttributeValue)

	assert.NoError(t, err)
	assert.NotNil(t, attributeValue)
}

func Test_convertRawAttributeValues_L_notList(t *testing.T) {
	expressionAttributeString := `{"L": "DogCow"}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_L_notListOfMaps(t *testing.T) {
	expressionAttributeString := `{"L": [ {"S": "Cookies"} , "DogCow", {"N": "3.14159"}]}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_M(t *testing.T) {
	expressionAttributeString := `{"M": {"Name": {"S": "Joe"}, "Age": {"N": "35"}}}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	attributeValue, err := convertRawAttributeValues(rawAttributeValue)

	assert.NoError(t, err)
	assert.NotNil(t, attributeValue)
}

func Test_convertRawAttributeValues_M_notMap(t *testing.T) {
	expressionAttributeString := `{"M": "DogCow"}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func Test_convertRawAttributeValues_M_notMapOfMaps(t *testing.T) {
	expressionAttributeString := `{"M": {"Name": {"S": "Joe"}, "Age": true}}`

	rawAttributeValue := make(map[string]interface{})
	err := json.Unmarshal([]byte(expressionAttributeString), &rawAttributeValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal the expression attribute value, %+v", err)
	}

	_, err = convertRawAttributeValues(rawAttributeValue)

	assert.Error(t, err)
}

func testChannelOfScanMethodReturnType() chan []map[string]types.AttributeValue {
	channel := make(chan []map[string]types.AttributeValue)

	go func() {
		channel <- testDynamoDbItems()
		close(channel)
	}()

	return channel
}
