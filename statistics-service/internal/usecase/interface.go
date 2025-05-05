package usecase

import (
	"statistics-service/internal/model"
)

type StatisticsRepository interface {
	SaveOrderEvent(event *model.OrderEvent) error
	SaveProductEvent(event *model.ProductEvent) error
	CountOrdersByUser(userID string) (int, error)
	OrdersGroupedByHour(userID string) (map[int]int, error)
	CountTotalUsers() (int, error)
	CountTotalProducts() (int, error)
}
type StatisticsService interface {
	SaveUserOrder(userID string, orderTime string) error
}
