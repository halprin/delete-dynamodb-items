package config

var tableName *string
var dynamoDbEndpoint *string
var filterExpression *string
var expressionAttributeNames *string
var expressionAttributeValues * string

func SetTableName(name string) {
	tableName = &name
}

func GetTableName() *string {
	return tableName
}

func SetDynamoDbEndpoint(endpoint string) {
	dynamoDbEndpoint = &endpoint
}

func GetDynamoDbEndpoint() *string {
	return dynamoDbEndpoint
}

func SetFilterExpression(expression string) {
	filterExpression = &expression
}

func GetFilterExpression() *string {
	return filterExpression
}

func SetExpressionAttributeNames(names string) {
	expressionAttributeNames = &names
}

func GetExpressionAttributeNames() *string {
	return expressionAttributeNames
}

func SetExpressionAttributeValues(values string) {
	expressionAttributeValues = &values
}

func GetExpressionAttributeValues() *string {
	return expressionAttributeValues
}
