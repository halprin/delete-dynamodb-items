package config

var tableName *string
var dynamoDbEndpoint *string

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
