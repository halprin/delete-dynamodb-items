#!/usr/bin/env bash
set -x

num_items="${1}"

table_name="mass-data"

aws dynamodb create-table --table-name "${table_name}" --attribute-definitions AttributeName=id,AttributeType=S --key-schema AttributeName=id,KeyType=HASH --billing-mode PAY_PER_REQUEST --endpoint-url http://127.0.0.1:8002 --no-cli-pager

items_preamble="{\"${table_name}\": ["
items_middle=""
items_ending=']}'
lorem_ipsum='Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec a efficitur nunc. Morbi fermentum sem metus, vel venenatis leo porttitor quis. Etiam maximus neque a pharetra viverra. Sed turpis lacus, blandit ac tortor elementum, scelerisque feugiat risus. Nam malesuada augue et purus aliquet, et semper dolor cursus. Suspendisse volutpat dolor nec efficitur rutrum. Aliquam leo libero, posuere eget vulputate in, luctus nec nibh. Donec eu tellus eu libero scelerisque molestie. Ut sed pretium nibh. Donec suscipit eget dui quis lacinia. Aliquam non pulvinar massa, nec blandit lectus. Cras sollicitudin rhoncus ex. Nunc ipsum dui, dictum in risus nec, convallis rutrum justo. In tempor dui nisl, in fringilla massa vehicula ac. Donec a ipsum luctus, venenatis magna ut, venenatis risus. Vivamus eu dapibus odio. Aenean dapibus urna orci, sed pharetra nunc dapibus ac. Praesent ornare, felis sit amet mattis faucibus, odio arcu laoreet arcu, eu blandit nisi turpis cursus enim.'

for ((index = 1 ; index <= num_items ; index++)); do
    rand_number=$((RANDOM % 100))  # not actually random, but lazy with the modulus
    current_request="{\"PutRequest\": {\"Item\": {\"id\": {\"S\": \"$(uuidgen)\"}, \"text\": {\"S\": \"${lorem_ipsum}\"}, \"number\": {\"N\": \"${rand_number}\"}}}}"
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
