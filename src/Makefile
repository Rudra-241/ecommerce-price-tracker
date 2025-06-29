SRC := cmd/main.go
COMPOSE := docker compose -f ../docker-compose.yml --env-file .env
all: build
build:
	@mkdir -p bin
	@go build -o bin/app $(SRC)
run: build
	bin/app
seed-selectors: build
	bin/app --job=seed-selectors
infra:
	$(COMPOSE) up -d $(SERVICES)
infra-auto:
	@services=""; \
	ss -Hltn 'sport = :5432' | grep -q . && echo "postgres: already on host, skipping" || services="postgres"; \
	ss -Hltn 'sport = :5672' | grep -q . && echo "rabbitmq: already on host, skipping" || services="$$services rabbitmq"; \
	if [ -z "$$services" ]; then \
		echo "Postgres and RabbitMQ already running on host; nothing to start."; \
	else \
		echo "Starting:$$services"; \
		$(COMPOSE) up -d $$services; \
	fi
infra-down:
	$(COMPOSE) down
clean:
	@rm -rf bin

.PHONY: all build run seed-selectors infra infra-auto infra-down clean
