.PHONY: postgres adminer migrate-up migrate-down

postgres:
	docker run --rm -ti -p 127.0.0.1:15432:5432 -e POSTGRES_PASSWORD=password postgres

adminer:
	docker run --rm -ti --network host adminer

migrate-up:
	go run cmd/migrate/migrate.go up

migrate-down:
	go run cmd/migrate/migrate.go down