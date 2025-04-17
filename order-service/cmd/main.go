package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net"
	"order-service/config"
	grpcserver "order-service/internal/adapter/grpc/server"
	"order-service/internal/adapter/handler"
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

	orderService := service.NewOrderService(repo)
	reviewService := service.NewReviewService(repo.(*postgres.PostgresRepo))

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, grpcserver.NewOrderGRPCServer(orderService, reviewService))

	go func() {
		fmt.Println("gRPC server started on port 50052")
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	r := gin.Default()
	handler.NewOrderHandler(r, orderService)
	r.Run(":8082")
}
