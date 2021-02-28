package main

import (
	"errors"
	"github.com/halprin/delete-dynamodb-items/dynamo"
	"github.com/halprin/delete-dynamodb-items/external/cli"
	"log"
	"os"
)

func main() {
	log.Println("Start")

	cli.FillConfig()

	tableName, err := getTableName()
	if err != nil {
		killExecution(err)
	}

	err = dynamo.DeleteAllItemsInTable(tableName)
	if err != nil {
		killExecution(err)
	}
	log.Println("Complete")
}

func getTableName() (string, error)  {
	if len(os.Args) < 2 {
		return "", errors.New("Provide a table name for the first argument")
	}
	return os.Args[1], nil
}

func killExecution(err error) {
	log.Println("Failure")
	log.Fatal(err.Error())
}
