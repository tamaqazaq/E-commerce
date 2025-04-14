package usecase

import "inventory-service/internal/model"

type DiscountUsecase interface {
	Create(discount *model.Discount) error
	GetProductsWithDiscount() ([]*model.Product, error)
	Delete(id string) error
}
