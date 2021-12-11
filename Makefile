compile:
	go build -o delete-dynamodb-items -v ./cmd/

runTestDynamoDB:
	docker-compose -f dynamodb-docker-compose.yml up -d
	echo "Navigate to http://localhost:8001 for DynamoDB-Admin"

loadTestData:
	./generate_mass_data.sh 500

test:
	go test ./...
