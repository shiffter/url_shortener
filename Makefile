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
gen-grpc:
	protoc -I ./api/proto ./api/proto/shortener/url_shortener.proto --go_out=./internal/proto_gen \
	--go_opt=paths=source_relative --go-grpc_out=./internal/proto_gen --go-grpc_opt=paths=source_relative
	#protoc -I ./api/proto --go_out=. --go-grpc_out=. --grpc-gateway_out=. ./api/proto/shortener/url_shortener.proto