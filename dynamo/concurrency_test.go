package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func Test_determineConcurrency_CalculatesOnDemandWithBillingMode(t *testing.T) {
	assert := assert.New(t)

	tableInfo := &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			BillingModeSummary: &dynamodb.BillingModeSummary{
				BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
			},
		},
	}

	concurrency, err := determineConcurrency(tableInfo)

	assert.Nil(err)
	assert.Equal(runtime.NumCPU(), concurrency)
}

func Test_determineConcurrency_CalculatesOnDemandWithProvisionBeingZero(t *testing.T) {
	assert := assert.New(t)

	tableInfo := &dynamodb.DescribeTableOutput{
		Table: &dynamodb.TableDescription{
			BillingModeSummary: nil,
			ProvisionedThroughput: &dynamodb.ProvisionedThroughputDescription{
				WriteCapacityUnits: aws.Int64(0),
			},
		},
	}

	concurrency, err := determineConcurrency(tableInfo)

	assert.Nil(err)
	assert.Equal(runtime.NumCPU(), concurrency)
}
