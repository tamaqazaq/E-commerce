package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Use(authMiddleware())

	inventoryURL, _ := url.Parse("http://localhost:8081")
	inventoryProxy := httputil.NewSingleHostReverseProxy(inventoryURL)

	orderURL, _ := url.Parse("http://localhost:8082")
	orderProxy := httputil.NewSingleHostReverseProxy(orderURL)

	r.Any("/inventory/*path", createProxyHandler(inventoryProxy, "/inventory"))
	r.Any("/orders/*path", createProxyHandler(orderProxy, "/orders"))

	r.Use(corsMiddleware())

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start API Gateway:", err)
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}
		c.Next()
	}
}

func createProxyHandler(proxy *httputil.ReverseProxy, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, prefix)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
