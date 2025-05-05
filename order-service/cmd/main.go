package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"net"
	"order-service/config"
	grpcserver "order-service/internal/adapter/grpc/server"
	"order-service/internal/adapter/handler"
	"order-service/internal/adapter/nats"
	"order-service/internal/adapter/postgres"
	"order-service/internal/app/service"
	pb "order-service/proto"
)

func main() {
	config.LoadEnv()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.GetEnv("DB_HOST"),
		config.GetEnv("DB_PORT"),
		config.GetEnv("DB_USER"),
		config.GetEnv("DB_PASSWORD"),
		config.GetEnv("DB_NAME"),
	)

	repo, err := postgres.NewPostgresRepository(dsn)
	if err != nil {
		panic(err)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Close()
	
	publisher := natsadapter.NewNatsPublisher(nc)

	orderService := service.NewOrderService(repo, publisher)
	reviewService := service.NewReviewService(repo.(*postgres.PostgresRepo))

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, grpcserver.NewOrderGRPCServer(orderService, reviewService, nc))

	go func() {
		listener, err := net.Listen("tcp", ":50052")
		if err != nil {
			panic(err)
		}
		fmt.Println("gRPC server started on port 50052")
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	r := gin.Default()
	handler.NewOrderHandler(r, orderService)
	r.Run(":8082")
}
