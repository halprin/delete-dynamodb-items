# delete-dynamodb-items

Bulk delete items from a DynamoDB table.

## Download

You can get `delete-dynamodb-items` from the [releases](https://github.com/halprin/delete-dynamodb-items/releases)
section of this GitHub repository.  There you will find downloads for your operating system and CPU architecture.

## Usage

_**Warning**: running this command will result in all the items in the specified table to be deleted immediately!  There
is no "are you sure?" prompt._

```shell
delete-dynamodb-items <table name> [--endpoint=URL] [--filter-expression=string] [--expression-attribute-names=JSON] [--expression-attribute-values=JSON]
```

The program uses the default AWS credential algorithm to determine what IAM entity and region is used.  E.g. the
`~/.aws/credentials` file, the `AWS_*` environment variables, etc.

### Filter Expressions

You can specify a special expression to filter out items you don't want deleted.  AKA, the item will be deleted if the
filter matches.  You can learn more about filter expressions in
[AWS's DynamoDB Developer Guide](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Scan.html#Scan.FilterExpression)
and the
[`filter-expression` section in the AWS CLI](https://docs.aws.amazon.com/cli/latest/reference/dynamodb/scan.html).

Use a combination of the `--filter-expression=`, `--expression-attribute-names=`, and `--expression-attribute-values=`
options.  These options work the same way as the options on the AWS CLI.

E.g. `--filter-expression='#k > :v' --expression-attribute-names='{"#k": "number"}' --expression-attribute-values='{":v": {"N": "50"}}'`

### Custom Endpoint

You can customize the DynamoDB endpoint with the `--endpoint=` (or `-e`) option.  Set it to the URL of the endpoint.
E.g. `--endpoint=http://localhost:8002`.  If unspecified, the default AWS endpoints are used.

## Contributing

Thank you for thinking of contributing!  Please see the [contributing guide](./CONTRIBUTING.md).

## Development

Run the following to compile your own copy from source.

```shell
make compile
```
