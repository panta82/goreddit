postgres:
	docker run --rm -ti --expose 127.0.0.1:15432:5432 -e POSTGRES_PASSWORD=password postgres

adminer:
	docker run --rm -ti --network host adminer

migrate:
	go run cmd/migrate/migrate.go up