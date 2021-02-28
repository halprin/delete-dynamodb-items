package main

import (
	"github.com/halprin/delete-dynamodb-items/dynamo"
	"github.com/halprin/delete-dynamodb-items/external/cli"
	"log"
)

func main() {
	log.Println("Start")

	cli.FillConfig()

	err := dynamo.DeleteAllItemsInTable()
	if err != nil {
		killExecution(err)
	}
	log.Println("Complete")
}

func killExecution(err error) {
	log.Println("Failure")
	log.Fatal(err.Error())
}
