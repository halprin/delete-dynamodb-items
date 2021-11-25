package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"math"
	"runtime"
)

func determineConcurrency(tableInfo *dynamodb.DescribeTableOutput) int {

	if isOnDemand(tableInfo) {
		concurrency := getOnDemandConcurrency()
		log.Printf("Given on demand concurrency, set concurrency to %d\n", concurrency)
		return concurrency
	}

	//the table has provisioned capacity

	numberOfItems := getNumberOfItems(tableInfo)
	tableSize := getTableSizeInBytes(tableInfo)

	roundedUpAverageItemSize := calculateAverageItemSize(tableSize, numberOfItems)
	totalBatchSize := float64(maxItemsPerBatchRequest) * roundedUpAverageItemSize

	writeCapacityUnits := getWriteCapacityUnits(tableInfo)
	rawConcurrency := float64(writeCapacityUnits) / totalBatchSize

	concurrency := 1
	if rawConcurrency > 1 {
		//possible truncation to size of int
		concurrency = int(rawConcurrency)
	}

	log.Printf("Given provisioned write capacity of %d, number of items %d, and table size %f KB, set concurrency to %d\n", writeCapacityUnits, numberOfItems, float64(tableSize) / float64(1024), concurrency)

	return concurrency
}

func isOnDemand(describeTable *dynamodb.DescribeTableOutput) bool {
	billingModeSummary := describeTable.Table.BillingModeSummary
	if billingModeSummary != nil {
		return *describeTable.Table.BillingModeSummary.BillingMode == dynamodb.BillingModePayPerRequest
	}

	return getWriteCapacityUnits(describeTable) == 0
}

func getWriteCapacityUnits(describeTable *dynamodb.DescribeTableOutput) int64 {
	return *describeTable.Table.ProvisionedThroughput.WriteCapacityUnits
}

func getNumberOfItems(describeTable *dynamodb.DescribeTableOutput) int64 {
	return *describeTable.Table.ItemCount
}

func getTableSizeInBytes(describeTable *dynamodb.DescribeTableOutput) int64 {
	return *describeTable.Table.TableSizeBytes
}

func getOnDemandConcurrency() int {
	//on demand's concurrency is the number of logical CPUs
	return runtime.NumCPU()
}

func calculateAverageItemSize(tableSize int64, numberOfItems int64) float64 {
	if tableSize == 0 || numberOfItems == 0 {
		//possible for the size or number of items to be 0 since they aren't always updated
		return 1.0
	}

	//truncates some of the digits if tableSize is too big
	//1024 to get Kilobytes
	averageItemSize := float64(tableSize) / float64(1024) / float64(numberOfItems)
	return math.Ceil(averageItemSize)
}
