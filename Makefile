.PHONY: help up down logs test build
help:
	@echo "targets: up down logs test build"
up:
	docker compose up -d --build
down:
	docker compose down
logs:
	docker compose logs --tail=100 whatsapp_go
test:
	cd src && go test ./...
build:
	docker compose build --pull whatsapp_go
