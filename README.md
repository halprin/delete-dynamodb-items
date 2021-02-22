# delete-dynamodb-items
Bulk delete items from a DynamoDB table.

## Download

You can get `delete-dynamodb-items` from the [releases](https://github.com/halprin/delete-dynamodb-items/releases)
section of this GitHub repository.  There you will find downloads for your operating system and CPU architecture.

## Usage

_**Warning**: running this command will result in all the items in the specified table to be deleted immediately!  There
is no "are you sure?" prompt._

```shell
delete-dynamodb-items <table name>
```

The program uses the default AWS credential algorithm to determine what IAM entity and region is used.  E.g. the
`~/.aws/credentials` file, the `AWS_*` environment variables, etc.

## Build

Run the following to compile your own copy from source.

```shell
go get -v -t -d ./cmd/
go build -o delete-dynamodb-items -v ./cmd/
```
