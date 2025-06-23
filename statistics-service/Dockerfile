FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o statistics-service ./cmd/main.go

FROM debian:bullseye-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/statistics-service .

EXPOSE 50053

CMD ["./statistics-service"]
