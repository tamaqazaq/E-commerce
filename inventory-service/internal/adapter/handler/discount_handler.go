package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"inventory-service/internal/model"
	"inventory-service/internal/usecase"
	"net/http"
)

type DiscountHandler struct {
	usecase usecase.DiscountUsecase
}

func NewDiscountHandler(router *gin.Engine, usecase usecase.DiscountUsecase) {
	h := &DiscountHandler{usecase: usecase}
	discount := router.Group("/discounts")
	{
		discount.POST("", h.CreateDiscount)
		discount.GET("/products", h.GetProductsWithDiscount)
		discount.DELETE("/:id", h.DeleteDiscount)
	}
}

func (h *DiscountHandler) CreateDiscount(c *gin.Context) {
	var discount model.Discount
	if err := c.ShouldBindJSON(&discount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	discount.ID = uuid.New().String()
	if err := h.usecase.Create(&discount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, discount)
}

func (h *DiscountHandler) GetProductsWithDiscount(c *gin.Context) {
	products, err := h.usecase.GetProductsWithDiscount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Discount deleted"})
}
