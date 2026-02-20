APP_SERVICE_NAME=app
APP_CONTAINER_NAME=mecanica_xpto
DB_SERVICE_NAME=db
DB_CONTAINER_NAME=db
APP_BINARY_PATH=/app/mecanica-xpto-api

.PHONY: init up down logs swag-generate-docker swag-run-docker test coverage coverage-html

init:
	cp .env-example .env
	docker-compose up -d --build

	@echo "Aguardando banco ficar pronto..."
	docker-compose exec $(DB_SERVICE_NAME) sh -c 'until pg_isready -U $$POSTGRES_USER -d $$POSTGRES_DB; do sleep 1; done'

	@echo "Aguardando container $(APP_CONTAINER_NAME) estar rodando..."
	@while [ -z "$$(docker-compose ps -q $(APP_SERVICE_NAME))" ] || \
		[ "$$(docker inspect -f '{{.State.Running}}' $$(docker-compose ps -q $(APP_SERVICE_NAME)))" != "true" ]; do \
		echo "Aguardando container..."; sleep 2; \
	done

	@echo "Container $(APP_CONTAINER_NAME) está rodando!"

	@echo "Aguardando app iniciar (sleep 10s dentro do container)..."
	docker-compose exec $(APP_SERVICE_NAME) sh -c 'sleep 10'

up:
	docker-compose up -d --build

down:
	docker-compose down

logs:
	docker-compose logs -f $(APP_SERVICE_NAME)

dev-up:
	cp .env-example .env
	docker-compose up -d dev

swag-generate: dev-up
	docker-compose exec dev sh -c "go install github.com/swaggo/swag/cmd/swag@latest && swag init -g internal/infrastructure/http/routes/routes.go --output ./docs --parseDependency --parseInternal"

test: dev-up
	docker-compose exec dev go test ./internal/usecase ./internal/adapter/http/handlers ./internal/adapter/http/routes ./cmd/api -v

coverage: dev-up
	docker-compose exec dev go test ./internal/usecase -coverprofile=usecasecov
	docker-compose exec dev go tool cover -func usecasecov

coverage-html: dev-up
	docker-compose exec dev go test ./internal/usecase -coverprofile=usecasecov
	docker cp $$(docker-compose ps -q dev):/app/usecasecov .
	go tool cover -html=usecasecov