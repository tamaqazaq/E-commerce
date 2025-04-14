package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"order-service/config"
	"order-service/internal/adapter/handler"
	"order-service/internal/adapter/postgres"
	"order-service/internal/app/service"
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

	r := gin.Default()
	usecase := service.NewOrderService(repo)
	handler.NewOrderHandler(r, usecase)
	r.Run(":8082")
}
