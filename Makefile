inmemory:
	docker-compose --profile inmemory up -d --build
postgres:
	docker-compose --profile postgres up -d --build
clean:
	docker-compose --profile postgres --profile inmemory down --volumes