package service

import (
	"statistics-service/internal/model"
	"statistics-service/internal/usecase"
	"time"
)

type StatisticsService struct {
	repo usecase.StatisticsRepository
}

func NewStatisticsService(repo usecase.StatisticsRepository) *StatisticsService {
	return &StatisticsService{repo: repo}
}

func (s *StatisticsService) GetUserOrderStats(userID string) (int, map[int]int, error) {
	total, err := s.repo.CountOrdersByUser(userID)
	if err != nil {
		return 0, nil, err
	}
	hourly, err := s.repo.OrdersGroupedByHour(userID)
	if err != nil {
		return 0, nil, err
	}
	return total, hourly, nil
}

func (s *StatisticsService) GetGeneralStats() (int, int, error) {
	users, err := s.repo.CountTotalUsers()
	if err != nil {
		return 0, 0, err
	}
	products, err := s.repo.CountTotalProducts()
	if err != nil {
		return users, 0, err
	}
	return users, products, nil
}
func (s *StatisticsService) SaveUserOrder(userID string, orderTime string) error {
	timestamp, err := time.Parse(time.RFC3339, orderTime)
	if err != nil {
		return err
	}

	event := &model.OrderEvent{
		UserID:    userID,
		OrderID:   "",
		Total:     0,
		Timestamp: timestamp,
	}

	return s.repo.SaveOrderEvent(event)
}
