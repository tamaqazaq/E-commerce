package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"inventory-service/config"
	"inventory-service/internal/adapter/handler"
	"inventory-service/internal/adapter/postgres"
	"inventory-service/internal/app/service"
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

	productRepo, err := postgres.NewPostgresProductRepository(dsn)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	productService := service.NewProductService(productRepo)
	handler.NewProductHandler(r, productService)

	discountRepo := postgres.NewPostgresDiscountRepository(productRepo.GetDB())
	discountService := service.NewDiscountService(discountRepo)
	handler.NewDiscountHandler(r, discountService)

	r.Run(":8081")
}
