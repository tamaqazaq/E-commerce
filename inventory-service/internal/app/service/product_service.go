package service

import (
	"inventory-service/internal/adapter/postgres"
	"inventory-service/internal/model"
	"inventory-service/internal/usecase"
)

type ProductService struct {
	repo postgres.ProductRepository
}

func NewProductService(repo postgres.ProductRepository) usecase.ProductUsecase {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(product *model.Product) error {
	return s.repo.Save(product)
}

func (s *ProductService) GetByID(id string) (*model.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) Update(id string, product *model.Product) error {
	product.ID = id
	return s.repo.Update(product)
}

func (s *ProductService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *ProductService) List() ([]*model.Product, error) {
	return s.repo.FindAll()
}
