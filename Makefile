# Variables
DOCKER_COMPOSE = docker-compose
IMAGE_NAME = order-list-golang-app
REST_API_IMAGE_NAME = order-list-golang-rest-api
DB_CONTAINER_NAME = order-list-golang-db-1
APP_CONTAINER_NAME = order-list-golang-app-1
REST_API_CONTAINER_NAME = order-list-golang-rest-api-1

# Default target
all: build up

# Build Docker images
build:
	$(DOCKER_COMPOSE) build

# Start Docker containers in detached mode
up:
	$(DOCKER_COMPOSE) up -d

# Stop and remove Docker containers
down:
	$(DOCKER_COMPOSE) down

# Run database migrations
migrate:
	$(DOCKER_COMPOSE) run --rm db psql -U user -d orders_db -f /migrations/create_orders_table.sql

# View logs from Docker containers
logs:
	$(DOCKER_COMPOSE) logs -f

# Open a shell in the app container
shell:
	docker exec -it $(APP_CONTAINER_NAME) /bin/sh

# Open a shell in the db container
dbshell:
	docker exec -it $(DB_CONTAINER_NAME) /bin/sh

# Rebuild Docker images and start containers
rebuild: down build up

# Clean up Docker images
clean:
	docker system prune -a --volumes -f

.PHONY: all build up down migrate logs shell dbshell rebuild clean
