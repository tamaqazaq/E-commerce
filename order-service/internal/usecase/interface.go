package usecase

import "order-service/internal/model"

type OrderUsecase interface {
	Create(order *model.Order) error
	GetByID(id string) (*model.Order, error)
	UpdateStatus(id, status string) error
	ListByUser(userID string) ([]*model.Order, error)
}
type ReviewUsecase interface {
	CreateReview(review *model.Review) error
	UpdateReview(review *model.Review) error
}
