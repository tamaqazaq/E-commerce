FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o inventory-service ./cmd/main.go

FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/inventory-service .

EXPOSE 8081 50051

CMD ["./inventory-service"]
