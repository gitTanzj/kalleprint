build: \
    go build -o bin/server cmd/main.go

test: \
    go test -v ./..

run: build \
    ./bin/server

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

seed:
	@for f in cmd/migrate/seeds/*.sql; do \
		echo "Running $$f..."; \
		docker compose exec -T db mysql -u root -proot kalleprint < $$f; \
	done
