package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_chunkItems(t *testing.T) {
	chunkedDynamoDbItems := chunkItems(testDynamoDbItems())

	for _, dynamoDbItems := range chunkedDynamoDbItems {
		assert.LessOrEqual(t, len(dynamoDbItems), maxItemsPerBatchRequest)
	}
}

func testDynamoDbItems() []map[string]*dynamodb.AttributeValue {
	sliceCapacity := 128

	sliceOfDynamoDbitems := make([]map[string]*dynamodb.AttributeValue, 0, sliceCapacity)

	for itemIndex := 0; itemIndex < sliceCapacity; itemIndex++ {
		sliceOfDynamoDbitems = append(sliceOfDynamoDbitems, map[string]*dynamodb.AttributeValue{
			fmt.Sprintf("dogcow%d", itemIndex): {
				S: aws.String(fmt.Sprintf("moof%d", itemIndex)),
			},
		})
	}

	return sliceOfDynamoDbitems
}