package service

import (
	"github.com/google/uuid"
	"order-service/internal/adapter/postgres"
	"order-service/internal/model"
	"order-service/internal/usecase"
	"time"
)

type ReviewService struct {
	repo postgres.ReviewRepository
}

func NewReviewService(repo *postgres.PostgresRepo) usecase.ReviewUsecase {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) CreateReview(review *model.Review) error {
	review.ID = uuid.New().String()
	now := time.Now()
	review.CreatedAt = now
	review.UpdatedAt = now
	return s.repo.Create(review)
}

func (s *ReviewService) UpdateReview(review *model.Review) error {
	review.UpdatedAt = time.Now()
	return s.repo.Update(review)
}
