inmemory:
	docker-compose --profile inmemory up --build
postgres:
	docker-compose --profile postgres up --build
clean:
	docker-compose --profile postgres --profile inmemory down --volumes
	rm testCoverage.out coverage.html
genmock:
	mockgen -source=internal/storage/interfaces.go -destination=internal/storage/mocks/mock_storage.go
test-coverage:
	go test -race ./... -tags unit -coverprofile=testCoverage.out
	go tool cover -html=testCoverage.out -o coverage.html
	firefox coverage.html    # On Linux
