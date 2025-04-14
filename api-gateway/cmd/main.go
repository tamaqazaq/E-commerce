package main

import (
	"api-gateway/proto"
	"api-gateway/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()

	client, err := service.NewGRPCClient("localhost:50051", "localhost:50052")
	if err != nil {
		log.Fatalf("could not initialize gRPC clients: %v", err)
	}

	r.POST("/products", func(c *gin.Context) {
		var productRequest proto.CreateProductRequest
		if err := c.ShouldBindJSON(&productRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		product, err := client.InventoryClient.CreateProduct(c, &productRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	r.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		product, err := client.InventoryClient.GetProductByID(c, &proto.GetProductRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, product)
	})
	r.PUT("/products/:id", func(c *gin.Context) {
		var productRequest proto.UpdateProductRequest
		if err := c.ShouldBindJSON(&productRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		product, err := client.InventoryClient.UpdateProduct(c, &productRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	r.DELETE("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		_, err := client.InventoryClient.DeleteProduct(c, &proto.DeleteProductRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
	})

	r.GET("/products", func(c *gin.Context) {
		products, err := client.InventoryClient.ListProducts(c, &proto.ListProductsRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	r.POST("/orders", func(c *gin.Context) {
		var orderRequest proto.CreateOrderRequest
		if err := c.ShouldBindJSON(&orderRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		order, err := client.OrderClient.CreateOrder(c, &orderRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	})

	r.GET("/orders/:id", func(c *gin.Context) {
		id := c.Param("id")
		order, err := client.OrderClient.GetOrderByID(c, &proto.GetOrderRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	})

	r.PUT("/orders/:id", func(c *gin.Context) {
		var orderStatusRequest proto.UpdateOrderStatusRequest
		if err := c.ShouldBindJSON(&orderStatusRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		order, err := client.OrderClient.UpdateOrderStatus(c, &orderStatusRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, order)
	})

	r.GET("/orders", func(c *gin.Context) {
		userID := c.DefaultQuery("user_id", "") // Assuming user_id is passed as query parameter
		orders, err := client.OrderClient.ListUserOrders(c, &proto.ListOrdersRequest{UserId: userID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	})

	r.Run(":8080")
}
