package usecase

import "inventory-service/internal/model"

type ProductUsecase interface {
	Create(product *model.Product) error
	GetByID(id string) (*model.Product, error)
	Update(id string, product *model.Product) error
	Delete(id string) error
	List() ([]*model.Product, error)
}
type ProductRepository interface {
	Save(product *model.Product) error
	FindByID(id string) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id string) error
	FindAll() ([]*model.Product, error)
}
