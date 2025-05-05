package usecase

import "inventory-service/internal/model"

type EventPublisher interface {
	PublishProductCreated(product *model.Product) error
	PublishProductUpdated(product *model.Product) error
	PublishProductDeleted(productID string) error
}
