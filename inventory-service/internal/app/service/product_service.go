package service

import (
	"github.com/google/uuid"
	"inventory-service/internal/model"
	"inventory-service/internal/usecase"
)

type ProductService struct {
	repo      usecase.ProductRepository
	cache     usecase.ProductCache
	publisher usecase.EventPublisher
}

func NewProductService(repo usecase.ProductRepository, cache usecase.ProductCache, publisher usecase.EventPublisher) usecase.ProductUsecase {
	return &ProductService{repo: repo, cache: cache, publisher: publisher}
}

func (s *ProductService) Create(product *model.Product) error {
	product.ID = uuid.New().String()
	if err := s.repo.Save(product); err != nil {
		return err
	}
	s.cache.Save(product)
	return s.publisher.PublishProductCreated(product)
}

func (s *ProductService) GetByID(id string) (*model.Product, error) {
	if p, ok := s.cache.GetByID(id); ok {
		return p, nil
	}
	return s.repo.FindByID(id)
}

func (s *ProductService) Update(id string, product *model.Product) error {
	product.ID = id
	if err := s.repo.Update(product); err != nil {
		return err
	}
	s.cache.Update(product)
	return s.publisher.PublishProductUpdated(product)
}

func (s *ProductService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.cache.Delete(id)
	return s.publisher.PublishProductDeleted(id)
}

func (s *ProductService) List() ([]*model.Product, error) {
	return s.cache.GetAll(), nil
}
