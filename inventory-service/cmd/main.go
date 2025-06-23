package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"inventory-service/config"
	"inventory-service/internal/adapter/cache"
	grpcserver "inventory-service/internal/adapter/grpc/server"
	"inventory-service/internal/adapter/handler"
	natsadapter "inventory-service/internal/adapter/nats"
	"inventory-service/internal/adapter/postgres"
	"inventory-service/internal/app/service"
	pb "inventory-service/proto"
	"net"
	"time"
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

	publisher := natsadapter.NewNatsPublisher(nc)

	memCache := cache.NewInMemoryProductCache()
	redisCache := cache.NewRedisProductCache()
	productCache := cache.NewMultiCache(memCache, redisCache)

	allProducts, err := repo.FindAll()
	if err != nil {
		panic(err)
	}
	productCache.LoadFromDB(allProducts)

	go func() {
		for {
			time.Sleep(12 * time.Hour)
			products, err := repo.FindAll()
			if err != nil {
				fmt.Println("Cache refresh failed:", err)
				continue
			}
			productCache.LoadFromDB(products)
			fmt.Println("Cache refreshed")
		}
	}()

	productService := service.NewProductService(repo, productCache, publisher)

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
	r.SetTrustedProxies(nil)
	handler.NewProductHandler(r, productService)
	r.Run(":8081")
}
