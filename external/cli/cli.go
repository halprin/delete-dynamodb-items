package cli

import (
	"github.com/halprin/delete-dynamodb-items/config"
	"github.com/teris-io/cli"
	"os"
)

func FillConfig() {
	endpointKey := "endpoint"
	tableNameCliArg := cli.NewArg("table name", "The name of the table for which all the items will be deleted").WithType(cli.TypeString)
	endpointCliOption := cli.NewOption(endpointKey, "A URL of the DynamoDB endpoint to use").WithChar('e').WithType(cli.TypeString)

	filterExpressionKey := "filter-expression"
	expressionAttributeNamesKey := "expression-attribute-names"
	expressionAttributeValuesKey := "expression-attribute-values"
	filterExpressionOption := cli.NewOption(filterExpressionKey, "A filter expression determines which items within the Scan results should be returned to you.  All of the other results are discarded.").WithType(cli.TypeString)
	expressionAttributeNamesOption := cli.NewOption(expressionAttributeNamesKey, "Way to specify names in the filter expression that are DynamoDB reserved words.").WithType(cli.TypeString)
	expressionAttributeValuesOption := cli.NewOption(expressionAttributeValuesKey, "Way to specify values in the filter expression that are DynamoDB reserved words.").WithType(cli.TypeString)

	parser := cli.New("Deletes all the items in a DynamoDB table").
		WithArg(tableNameCliArg).
		WithOption(endpointCliOption).
		WithOption(filterExpressionOption).
		WithOption(expressionAttributeNamesOption).
		WithOption(expressionAttributeValuesOption)

	invocation, arguments, options, err := parser.Parse(os.Args)
	help, helpExistsInOptions := options["help"]

	if err != nil {
		_ = parser.Usage(invocation, os.Stdout)
		os.Exit(1)
	} else if helpExistsInOptions && help == "true" {
		_ = parser.Usage(invocation, os.Stdout)
		os.Exit(0)
	}

	tableName := arguments[0]
	config.SetTableName(tableName)

	endpoint, endpointExistsInOptions := options[endpointKey]
	if endpointExistsInOptions {
		config.SetDynamoDbEndpoint(endpoint)
	}

	expression, expressionExistsInOptions := options[filterExpressionKey]
	if expressionExistsInOptions {
		config.SetFilterExpression(expression)
	}

	names, namesExistsInOptions := options[expressionAttributeNamesKey]
	if namesExistsInOptions {
		config.SetExpressionAttributeNames(names)
	}

	values, valuesExistsInOptions := options[expressionAttributeValuesKey]
	if valuesExistsInOptions {
		config.SetDynamoDbEndpoint(values)
	}
}
