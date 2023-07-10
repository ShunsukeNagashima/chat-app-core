.PHONY: up down test logs migrate lint gen-mocks

up:
		docker compose up -d

down:
		docker compose down

test: ## Execute test
		go test -race -shuffle=on ./...

logs: ## Tail docker compose logs
		docker compose logs -f

migrate:
		go run ./scripts/migration.go

lint:
		golangci-lint run --config=./.golangci.yml

gen-mocks:
		go generate ./pkg/domain/...
