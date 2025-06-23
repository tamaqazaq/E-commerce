package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

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
	config.LoadEnv()

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

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("NATS error: ", err)
	}
	defer nc.Close()

	natslistener.SubscribeToEvents(nc, repo)

	go func() {
		for {
			time.Sleep(10 * time.Second)

			types := []string{"order", "item"}
			messageType := types[rand.Intn(len(types))]

			data := map[string]interface{}{
				"event":  "hourly_stats_update",
				"time":   time.Now().Format(time.RFC3339),
				"source": "statistics-service",
				"type":   messageType,
			}
			payload, err := json.Marshal(data)
			if err != nil {
				log.Println("Failed to marshal hourly stats update:", err)
				continue
			}
			err = nc.Publish("ap2.statistics.event.updated", payload)
			if err != nil {
				log.Println("Failed to publish hourly stats update:", err)
			} else {
				log.Printf("Published hourly stats update (%s) to ap2.statistics.event.updated\n", messageType)
			}
		}
	}()

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
