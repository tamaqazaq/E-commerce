package usecase

import "order-service/internal/model"

type EventPublisher interface {
	PublishOrderCreated(order *model.Order) error
	PublishOrderUpdated(order *model.Order) error
}
