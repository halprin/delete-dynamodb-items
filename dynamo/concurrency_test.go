package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func Test_determineConcurrency_CalculatesOnDemandWithBillingMode(t *testing.T) {
	assert := assert.New(t)

	tableInfo := &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			BillingModeSummary: &types.BillingModeSummary{
				BillingMode: types.BillingModePayPerRequest,
			},
		},
	}

	concurrency := determineConcurrency(tableInfo)

	assert.Equal(runtime.NumCPU(), concurrency)
}

func Test_determineConcurrency_CalculatesOnDemandWithProvisionBeingZero(t *testing.T) {
	assert := assert.New(t)

	tableInfo := &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			BillingModeSummary: nil,
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				WriteCapacityUnits: aws.Int64(0),
			},
		},
	}

	concurrency := determineConcurrency(tableInfo)

	assert.Equal(runtime.NumCPU(), concurrency)
}

func Test_determineConcurrency_CalculatesProvision(t *testing.T) {
	assert := assert.New(t)

	tableInfo := &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				WriteCapacityUnits: aws.Int64(2631),
			},
			ItemCount:      aws.Int64(36001),
			TableSizeBytes: aws.Int64(506060108),
		},
	}

	concurrency := determineConcurrency(tableInfo)

	//average item size in KB = ceiling(table size in bytes / 1024 / number of items)
	//average batch size = max items per batch (25) * average item size in KB
	//raw concurrency = floor(write capacity units / average batch size)

	//So...
	//average item size in KB = 50600108 / 1024 / 36001 = 13.727377690029444 = 14
	//average batch size = 25 * 14 = 350
	//raw concurrency = 2631 / 350 = 7.517142857142857 = 7
	assert.Equal(7, concurrency)
}

func Test_determineConcurrency_CalculatesProvisionWithTableSizeZero(t *testing.T) {
	assert := assert.New(t)

	tableInfo := &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				WriteCapacityUnits: aws.Int64(2631),
			},
			ItemCount:      aws.Int64(36001),
			TableSizeBytes: aws.Int64(0),
		},
	}

	concurrency := determineConcurrency(tableInfo)

	//So...
	//average item size in KB = 1 (because of 0 table size)
	//average batch size = 25 * 1 = 25
	//raw concurrency = 2631 / 25 = 105.24 = 105
	assert.Equal(105, concurrency)
}

func Test_determineConcurrency_CalculatesProvisionWithItemCountZero(t *testing.T) {
	assert := assert.New(t)

	tableInfo := &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				WriteCapacityUnits: aws.Int64(2631),
			},
			ItemCount:      aws.Int64(0),
			TableSizeBytes: aws.Int64(506060108),
		},
	}

	concurrency := determineConcurrency(tableInfo)

	//average item size in KB = ceiling(table size in bytes / 1024 / number of items)
	//average batch size = max items per batch (25) * average item size in KB
	//raw concurrency = floor(write capacity units / average batch size)

	//So...
	//average item size in KB = 1 (because of 0 table size)
	//average batch size = 25 * 1 = 25
	//raw concurrency = 2631 / 25 = 105.24 = 105
	assert.Equal(105, concurrency)
}

func Test_determineConcurrency_CalculatesProvisionWithConcurencyLessThanOne(t *testing.T) {
	assert := assert.New(t)

	tableInfo := &dynamodb.DescribeTableOutput{
		Table: &types.TableDescription{
			ProvisionedThroughput: &types.ProvisionedThroughputDescription{
				WriteCapacityUnits: aws.Int64(34),
			},
			ItemCount:      aws.Int64(36001),
			TableSizeBytes: aws.Int64(506060108),
		},
	}

	concurrency := determineConcurrency(tableInfo)

	//average item size in KB = ceiling(table size in bytes / 1024 / number of items)
	//average batch size = max items per batch (25) * average item size in KB
	//raw concurrency = floor(write capacity units / average batch size)

	//So...
	//average item size in KB = 50600108 / 1024 / 36001 = 13.727377690029444 = 14
	//average batch size = 25 * 14 = 350
	//raw concurrency = 34 / 350 = 0.097142857142857 = 1
	assert.Equal(1, concurrency)
}
