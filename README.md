
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

## Prerequisites

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

## Using the Makefile

The Makefile included in this project provides a streamlined way to manage the build, deployment, testing, and cleanup processes for the Order Listing Service. It abstracts away complex Docker and database commands, making it easier to work with the application's various components, including REST API, GraphQL, and gRPC services.

### Prerequisites

Before you begin using the Makefile, ensure you have the following tools installed on your machine:

- Docker
- Docker Compose
- Go
- grpcurl (for testing gRPC services)
- Curl (for testing REST and GraphQL endpoints)

### Commands Overview

The Makefile supports several commands that you can use to control different aspects of the application. Below is a detailed list of available commands and their purposes:

| Command           | Description                                                                                             |
|-------------------|---------------------------------------------------------------------------------------------------------|
| `make all`        | Builds the Docker images and starts all containers in detached mode.                                     |
| `make build`      | Builds the Docker images for the application.                                                            |
| `make up`         | Starts all Docker containers in detached mode.                                                           |
| `make down`       | Stops and removes all Docker containers.                                                                 |
| `make migrate`    | Runs database migrations to set up the required tables in PostgreSQL.                                    |
| `make seed`       | Inserts sample data into the PostgreSQL database for testing purposes.                                   |
| `make logs`       | Displays logs from all running Docker containers.                                                        |
| `make shell`      | Opens a shell session inside the application container.                                                  |
| `make dbshell`    | Opens a shell session inside the database container.                                                     |
| `make rebuild`    | Stops, removes, builds, and starts the Docker containers again.                                          |
| `make clean`      | Removes all Docker images, containers, networks, and volumes to free up space.                           |
| `make test-rest`  | Tests the REST API endpoint using curl.                                                                  |
| `make test-graphql` | Tests the GraphQL API using curl.                                                                      |
| `make test-grpc`  | Tests the gRPC service using grpcurl.                                                                    |

### Running the Application

To start the application using the Makefile, follow these steps:

1. **Build and Start Containers**

   Execute the following command to build the Docker images and start the containers:

   ```bash
   make all
   ```

   This command is equivalent to running `docker-compose build` followed by `docker-compose up -d`.

2. **Run Migrations**

   After the containers are up and running, set up the necessary tables in the database by executing:

   ```bash
   make migrate
   ```

   This command runs the SQL script located at `/migrations/create_orders_table.sql` inside the database container.

3. **Insert Sample Data**

   To populate the database with sample data for testing, run:

   ```bash
   make seed
   ```

   This command inserts predefined data entries into the `orders` table.

4. **Testing the Services**

   - **REST API**

     Test the REST API endpoint using the following command:

     ```bash
     make test-rest
     ```

     This will send a GET request to `http://localhost:8080/order` and display the response.

   - **GraphQL**

     Test the GraphQL API by executing:

     ```bash
     make test-graphql
     ```

     This command sends a query to the GraphQL server available at `http://localhost:8081/playground`.

   - **gRPC**

     Use the following command to test the gRPC service:

     ```bash
     make test-grpc
     ```

     This command will call the `ListOrders` method on the gRPC server running at `localhost:50051`.

5. **Clean Up**

   When you're done with the application and want to free up resources, run:

   ```bash
   make clean
   ```

   This command removes all Docker images, containers, networks, and volumes associated with the project.

### Additional Commands

- **Access Shell in Containers**

  - To open a shell in the application container, use:

    ```bash
    make shell
    ```

  - To open a shell in the database container, use:

    ```bash
    make dbshell
    ```

- **View Logs**

  To view live logs from all running containers, execute:

  ```bash
  make logs
  ```

This section of the README should now guide users effectively on how to use the Makefile to manage the application lifecycle. Make sure to customize any part of this documentation to better fit the specific needs or changes in your project setup.