# E-commerce Microservices Platform

A scalable e-commerce platform built using a microservices architecture in **Go**. The system includes product management, order processing, layered caching, statistics, and asynchronous communication using **gRPC** and **NATS** — all containerized with Docker.

## Features

- Microservices with gRPC communication  
- RESTful API via centralized API Gateway (Gin)  
- PostgreSQL for persistent data  
- NATS for asynchronous pub/sub messaging  
- In-Memory + Redis-style layered caching  
- Statistics tracking  
- Fully Dockerized setup


## Architecture

```

\[ Client ] ⇄ \[ API Gateway (REST) ] ⇄ \[ Services (gRPC) ]
⇵
\[ NATS Broker ]

````

All client requests go through the API Gateway, which communicates with internal services using gRPC. Services publish and consume events via NATS for asynchronous communication.

## Tech Stack

- **Language:** Go 1.23.4
- **Framework:** Gin (for REST API)  
- **RPC Protocol:** gRPC  
- **Database:** PostgreSQL  
- **Messaging:** NATS  
- **Caching:** In-Memory + Redis-style pattern  
- **Containerization:** Docker + Docker Compose  

## Getting Started

### 1. Clone the Repository (use `cache` branch)
```bash
git clone --branch cache https://github.com/tamaqazaq/E-commerce.git
cd E-commerce
````

### 2. Launch the System

```bash
docker-compose up --build
```

### 3. Available Endpoints

* `POST /products` – create a product
* `GET /products/:id` – get product by ID 
* `GET /products` - Get list of all products
* `PUT /products/:id` - Update a product
* `DELETE /products/:id` - Delete a product
* `POST /orders` - Create a new order
* `GET /orders` - Get all orders (optionally filter by user_id via query param)
* `GET /orders/:id.` - get orders by user
* `PUT /orders/:id` - Update order status
* `GET /stats/user-orders/:user_id` - Get statistics about orders for a specific user
* `GET /stats/user` - Get general user statistics

## Project Structure

```
E-commerce/
├── api-gateway/
├── inventory-service/
├── order-service/
├── statistics-service/
├── docker-compose.yml
└── README.md
```

Each service follows clean architecture:

* `cmd/main.go` — entry point
* `internal/` — layered domain (handler, service, repository)
* `proto/` — gRPC service definitions

## Caching Strategy

The `inventory-service` uses a multi-layered cache:

1. In-memory cache for fast access
2. (Optional) Redis-style external cache
3. Automatic refresh from DB every 12 hours

## Event-driven Communication

* Services communicate using NATS topics
* For example, `order-service` publishes order events
* `statistics-service` subscribes and updates analytics

## Planned coverage:

* Core business logic (services)
* Repository layer (PostgreSQL/cache)
* gRPC and HTTP handlers (mock-based)


## Development Status

* Core services implemented and connected
* gRPC + HTTP API working
* Dockerized microservice architecture

## License

This project is licensed under the **MIT License**.
Feel free to use, modify, and share with proper attribution.

