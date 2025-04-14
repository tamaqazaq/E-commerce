package service

import (
	"inventory-service/internal/adapter/postgres"
	"inventory-service/internal/model"
	"inventory-service/internal/usecase"
)

type DiscountService struct {
	repo postgres.DiscountRepository
}

func NewDiscountService(repo postgres.DiscountRepository) usecase.DiscountUsecase {
	return &DiscountService{repo: repo}
}

func (s *DiscountService) Create(discount *model.Discount) error {
	return s.repo.Save(discount)
}

func (s *DiscountService) GetProductsWithDiscount() ([]*model.Product, error) {
	return s.repo.GetProductsWithDiscount()
}

func (s *DiscountService) Delete(id string) error {
	return s.repo.Delete(id)
}
