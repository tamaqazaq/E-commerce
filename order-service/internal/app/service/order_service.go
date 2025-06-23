package service

import (
	"github.com/google/uuid"
	"order-service/internal/model"
	"order-service/internal/usecase"
	"time"
)

type OrderService struct {
	repo      usecase.OrderRepository
	publisher usecase.EventPublisher
	cache     usecase.OrderCache
}

func NewOrderService(repo usecase.OrderRepository, publisher usecase.EventPublisher, cache usecase.OrderCache) usecase.OrderUsecase {
	return &OrderService{
		repo:      repo,
		publisher: publisher,
		cache:     cache,
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
	order.Timestamp = time.Now()

	if err := s.repo.Save(order); err != nil {
		return err
	}
	//return fmt.Errorf("db is disabled: cannot save order")
	s.cache.Save(order)
	return s.publisher.PublishOrderCreated(order)

}

func (s *OrderService) GetByID(id string) (*model.Order, error) {
	if o, ok := s.cache.GetByID(id); ok {
		return o, nil
	}
	return s.repo.FindByID(id)
	//return nil, fmt.Errorf("not found in cache")
}

func (s *OrderService) UpdateStatus(id, status string) error {
	if err := s.repo.UpdateStatus(id, status); err != nil {
		return err
	}
	order, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	s.cache.Update(order)
	return s.publisher.PublishOrderUpdated(order)
}

func (s *OrderService) ListByUser(userID string) ([]*model.Order, error) {
	if list, ok := s.cache.ListByUser(userID); ok {
		return list, nil
	}
	return s.repo.FindByUserID(userID)
	//return nil, fmt.Errorf("not found in cache")
}
