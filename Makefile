# Variables
APP_CONTAINER_NAME=order-list-golang-app-1
REST_API_CONTAINER_NAME=order-list-golang-rest-api-1
DB_CONTAINER_NAME=order-list-golang-db-1
DOCKER_COMPOSE=docker-compose

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
	$(DOCKER_COMPOSE) exec db psql -U $$DB_USER -d $$DB_NAME -f /migrations/create_orders_table.sql

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

# Insert sample data into the database
seed:
	$(DOCKER_COMPOSE) exec db psql -U $$DB_USER -d $$DB_NAME -c "INSERT INTO orders (user_id, product_id, quantity, status, created_at, updated_at) VALUES (1, 101, 2, 'Pending', '2024-07-30 10:00:00', '2024-07-30 10:00:00'), (2, 102, 1, 'Shipped', '2024-07-30 11:00:00', '2024-07-30 11:00:00'), (3, 103, 5, 'Delivered', '2024-07-30 12:00:00', '2024-07-30 12:00:00');"

# Test REST API using curl
test-rest:
	curl http://localhost:8080/order

# Test GraphQL API using curl
test-graphql:
	curl -X POST http://localhost:8081/query -H "Content-Type: application/json" -d '{"query": "{ listOrders { id user_id product_id quantity status created_at updated_at } }"}'

# Test gRPC service using grpcurl
test-grpc:
	grpcurl -plaintext -d '{}' localhost:50051 order.OrderService/ListOrders

.PHONY: all build up down migrate logs shell dbshell rebuild clean seed test-rest test-graphql test-grpc

