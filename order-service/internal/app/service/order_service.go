package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"order-service/internal/adapter/postgres"
	"order-service/internal/model"
	"order-service/internal/usecase"
	"time"
)

type OrderService struct {
	repo postgres.OrderRepository
	nc   *nats.Conn
}

func NewOrderService(repo postgres.OrderRepository, nc *nats.Conn) usecase.OrderUsecase {
	return &OrderService{
		repo: repo,
		nc:   nc,
	}
}

func (s *OrderService) Create(order *model.Order) error {
	total := 0.0
	order.ID = uuid.New().String()
	for i := range order.Items {
		order.Items[i].ID = uuid.New().String()
		order.Items[i].OrderID = order.ID
		total += order.Items[i].Price * float64(order.Items[i].Quantity)
	}
	order.Total = total
	order.Status = "pending"

	if err := s.repo.Save(order); err != nil {
		return err
	}

	event := map[string]interface{}{
		"action": "order.created",
		"time":   time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"order_id": order.ID,
			"user_id":  order.UserID,
			"total":    order.Total,
		},
	}
	payload, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error marshaling order.created event:", err)
		return nil
	}
	if err := s.nc.Publish("order.created", payload); err != nil {
		fmt.Println("Error publishing order.created event:", err)
	}

	return nil
}

func (s *OrderService) GetByID(id string) (*model.Order, error) {
	return s.repo.FindByID(id)
}

func (s *OrderService) UpdateStatus(id, status string) error {
	return s.repo.UpdateStatus(id, status)
}

func (s *OrderService) ListByUser(userID string) ([]*model.Order, error) {
	return s.repo.FindByUserID(userID)
}
