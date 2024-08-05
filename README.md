
# Order Listing Service

This project is an implementation of a simple order listing service using gRPC, GraphQL, and REST API. It demonstrates how to create a multi-service application with a shared PostgreSQL database using Docker Compose.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Create .env file](#create-env-file)
3. [Run the Application](#run-the-application)
4. [Run Migrations](#run-migrations)
5. [Insert Sample Data](#insert-sample-data)
6. [Testing the Services](#testing-the-services)
   - [REST API](#rest-api)
   - [GraphQL](#graphql)
   - [gRPC](#grpc)
7. [Clean Up](#clean-up)

## 1. Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/doc/install)
- [grpcurl](https://github.com/fullstorydev/grpcurl#installation) (for testing gRPC services)

## 2. Create .env file

Create a `.env` file in the root directory of the project with the following content:

```ini
DB_USER=user
DB_PASSWORD=password
DB_NAME=orders_db
DB_HOST=db
DB_PORT=5432
```

## 3. Run the Application

First, make sure Docker is running, then execute the following commands:

```bash
docker-compose build
docker-compose up -d
```

This will build the application and start all the services defined in the `docker-compose.yml` file.

## 4. Run Migrations

To create the required tables in the PostgreSQL database, execute the following command:

```bash
 docker-compose exec db sh 
```

```bash
psql -U user -d orders_db
```

```bash
  CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


  INSERT INTO orders (user_id, product_id, quantity, status, created_at, updated_at) VALUES
    (1, 101, 2, 'Pending', '2024-07-30 10:00:00', '2024-07-30 10:00:00'),
    (2, 102, 1, 'Shipped', '2024-07-30 11:00:00', '2024-07-30 11:00:00'),
    (3, 103, 5, 'Delivered', '2024-07-30 12:00:00', '2024-07-30 12:00:00');
```

```bash
  exit
```

## 5. Testing the Services

### REST API

You can test the REST API using `curl` or any API testing tool like Postman. The REST API will be available at `http://localhost:8080`.

```bash
curl http://localhost:8080/order
```

### GraphQL

You can test the GraphQL API using the GraphQL Playground available at `http://localhost:8081/playground`.

Example Query:

```graphql
query {
  listOrders {
    id
    user_id
    product_id
    quantity
    status
    created_at
    updated_at
  }
}
```

### gRPC

You can test the gRPC service using `grpcurl`. The gRPC service will be available at `localhost:50051`.

```bash
grpcurl -plaintext localhost:50051 list
```

To call the `ListOrders` method:

```bash
grpcurl -plaintext -d '{}' localhost:50051 order.OrderService/ListOrders
```

## 6. Clean Up

To stop and remove all the containers, network, and volumes, run the following command:

```bash
docker-compose down -v
```


