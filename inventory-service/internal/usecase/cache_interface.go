package usecase

import "inventory-service/internal/model"

type ProductCache interface {
	GetAll() []*model.Product
	GetByID(id string) (*model.Product, bool)
	Save(product *model.Product)
	Update(product *model.Product)
	Delete(id string)
	LoadFromDB(products []*model.Product)
}
