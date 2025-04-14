package service

import (
	"github.com/google/uuid"
	"order-service/internal/adapter/postgres"
	"order-service/internal/model"
	"order-service/internal/usecase"
)

type OrderService struct {
	repo postgres.OrderRepository
}

func NewOrderService(repo postgres.OrderRepository) usecase.OrderUsecase {
	return &OrderService{repo: repo}
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
	return s.repo.Save(order)
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
