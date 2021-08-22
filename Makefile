compile:
	go build -o delete-dynamodb-items -v ./cmd/

runTestDynamoDB:
	docker-compose -f dynamodb-docker-compose.yml up -d

loadTestData:
	./generate_mass_data.sh 500
