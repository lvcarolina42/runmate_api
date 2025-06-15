deps:
	@docker compose up -d

deps-down:
	@docker compose down -v --remove-orphans

run: deps
	@sleep 2
	@go run cmd/main.go

.PHONY: deps runs