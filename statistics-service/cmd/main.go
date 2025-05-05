package main

import (
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"statistics-service/config"
	"statistics-service/internal/adapter/grpc/server"
	"statistics-service/internal/adapter/natslistener"
	"statistics-service/internal/adapter/postgres"
	"statistics-service/internal/app/service"
	pb "statistics-service/proto"
)

func main() {
	// Load env
	config.LoadEnv()

	// DB init
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.GetEnv("DB_HOST"),
		config.GetEnv("DB_PORT"),
		config.GetEnv("DB_USER"),
		config.GetEnv("DB_PASSWORD"),
		config.GetEnv("DB_NAME"),
	)
	repo, err := postgres.NewStatisticsRepository(dsn)
	if err != nil {
		log.Fatal("DB error: ", err)
	}

	// NATS connect
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("NATS error: ", err)
	}
	defer nc.Close()

	// Subscribe to NATS topics
	natslistener.SubscribeToEvents(nc, repo)

	// gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterStatisticsServiceServer(grpcServer, server.NewStatisticsGRPCServer(service.NewStatisticsService(repo)))

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal("gRPC listen error: ", err)
	}
	log.Println("Statistics gRPC server running on :50053")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("gRPC server error: ", err)
	}
}
