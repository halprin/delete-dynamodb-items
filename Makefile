compile:
	go build -o delete-dynamodb-items -v ./cmd/

runTestDynamoDB:
	docker-compose -f dynamodb-docker-compose.yml up -d
	echo "Navigate to http://localhost:8001 for DynamoDB-Admin"

loadTestData:
	./generate_mass_data.sh 500

test: unitTest integrationTest

unitTest:
	go test ./...

integrationTest: compile runTestDynamoDB
	./generate_mass_data.sh 1000
	AWS_REGION=us-east-1 AWS_ACCESS_KEY_ID=DogCow AWS_SECRET_ACCESS_KEY=Moof ./delete-dynamodb-items mass-data --endpoint=http://127.0.0.1:8002
	aws dynamodb describe-table --table-name mass-data --endpoint-url http://127.0.0.1:8002 | jq --exit-status '.Table.ItemCount == 0'
	docker-compose -f dynamodb-docker-compose.yml stop
