#!/usr/bin/env bash

num_items="${1}"

table_name="mass-data"

aws dynamodb create-table --table-name "${table_name}" --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --billing-mode PAY_PER_REQUEST --endpoint-url http://127.0.0.1:8002

items_preamble="{\"${table_name}\": ["
items_middle=""
items_ending=']}'

for ((index = 1 ; index <= num_items ; index++)); do
    current_request="{\"PutRequest\": {\"Item\": {\"id\": {\"S\": \"$(uuidgen)\"}}}}"
    items_middle="${items_middle}${current_request},"
    if [[ $((index % 25)) == 0 ]]; then
        items_middle=${items_middle::${#items_middle}-1}
        aws dynamodb batch-write-item --request-items "${items_preamble}${items_middle}${items_ending}" --endpoint-url http://127.0.0.1:8002
        items_middle=""
    fi
done

if [[ -n "${items_middle}" ]]; then
    items_middle=${items_middle::${#items_middle}-1}
    aws dynamodb batch-write-item --request-items "${items_preamble}${items_middle}${items_ending}" --endpoint-url http://127.0.0.1:8002
fi
