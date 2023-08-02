.PHONY: up down test logs migrate lint gen-mocks

up:
		@docker-compose up -d
		@sleep 5 \
		&& export AWS_ACCESS_KEY_ID=dummy \
            AWS_SECRET_ACCESS_KEY=dummy \
            AWS_SESSION_TOKEN=dummy \
            AWS_DEFAULT_REGION=ap-northeast-1 \
		&& aws --endpoint-url=http://localhost:4566 secretsmanager create-secret \
				--name "app/local/AppSecrets" \
				--secret-string "$(cat ./secrets/firebase-credentials.json)"

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

build: ## Build docker image to deploy
		docker build -t chat-app-core:latest \
						--target deploy .
