package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"inventory-service/config"
	grpcserver "inventory-service/internal/adapter/grpc/server"
	"inventory-service/internal/adapter/handler"
	"inventory-service/internal/adapter/postgres"
	"inventory-service/internal/app/service"
	pb "inventory-service/proto"
	"net"
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

	repo, err := postgres.NewPostgresProductRepository(dsn)
	if err != nil {
		panic(err)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	productService := service.NewProductService(repo, nc)

	go func() {
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer()
		pb.RegisterInventoryServiceServer(grpcServer, grpcserver.NewInventoryGRPCServer(productService))
		fmt.Println("gRPC server started on port 50051")
		if err := grpcServer.Serve(listener); err != nil {
			panic(err)
		}
	}()

	r := gin.Default()
	handler.NewProductHandler(r, productService)
	r.Run(":8081")
}
