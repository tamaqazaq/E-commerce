package usecase

import "order-service/internal/model"

type OrderCache interface {
	GetByID(id string) (*model.Order, bool)
	ListByUser(userID string) ([]*model.Order, bool)
	Save(order *model.Order)
	Update(order *model.Order)
	Delete(id string)
	LoadFromDB(orders []*model.Order)
}
