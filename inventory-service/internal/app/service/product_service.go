package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"inventory-service/internal/adapter/postgres"
	"inventory-service/internal/model"
	"inventory-service/internal/usecase"
	"time"
)

type ProductService struct {
	repo postgres.ProductRepository
	nc   *nats.Conn
}

func NewProductService(repo postgres.ProductRepository, nc *nats.Conn) usecase.ProductUsecase {
	return &ProductService{
		repo: repo,
		nc:   nc,
	}
}

func (s *ProductService) Create(product *model.Product) error {
	product.ID = uuid.New().String()
	err := s.repo.Save(product)
	if err != nil {
		return err
	}
	s.publishEvent("product.created", product)
	return nil
}

func (s *ProductService) GetByID(id string) (*model.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) Update(id string, product *model.Product) error {
	product.ID = id
	err := s.repo.Update(product)
	if err != nil {
		return err
	}
	s.publishEvent("product.updated", product)
	return nil
}

func (s *ProductService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	s.publishEvent("product.deleted", &model.Product{ID: id})
	return nil
}

func (s *ProductService) List() ([]*model.Product, error) {
	return s.repo.FindAll()
}

func (s *ProductService) publishEvent(topic string, product *model.Product) {
	event := map[string]interface{}{
		"action": topic,
		"time":   time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"id":       product.ID,
			"name":     product.Name,
			"category": product.Category,
			"price":    product.Price,
			"stock":    product.Stock,
		},
	}
	payload, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error marshaling event:", err)
		return
	}

	if err := s.nc.Publish(topic, payload); err != nil {
		fmt.Println("Error publishing to NATS:", err)
	}
}
