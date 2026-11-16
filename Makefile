.PHONY: help up down logs test build backup-storages gowa-up
help:
	@echo "targets: up down logs test build backup-storages gowa-up (on servers use gowa-up, not up)"
up:
	docker compose up -d --build
# Copies the WhatsApp session folder (mount of /app/storages) to <folder>.bak-<time>; read-only on the live session.
# Source = mount of the running container, else ./storages. A failed copy stops gowa-up.
backup-storages:
	@src=$$(docker inspect -f '{{range .Mounts}}{{if eq .Destination "/app/storages"}}{{.Source}}{{end}}{{end}}' $$(docker compose ps -q whatsapp_go) 2>/dev/null); \
	src=$${src:-$(CURDIR)/storages}; \
	test -d "$$src" || { echo "backup-storages: $$src not found, nothing to back up"; exit 0; }; \
	dst="$$src.bak-$$(date +%Y%m%d-%H%M%S)"; \
	cp -a "$$src" "$$dst" && echo "backup-storages: $$src -> $$dst ($$(du -sh $$dst | cut -f1))"
# Safe recreate: back up the session first, then up.
gowa-up: backup-storages up
down:
	docker compose down
logs:
	docker compose logs --tail=100 whatsapp_go
test:
	cd src && go test ./...
build:
	docker compose build --pull whatsapp_go
