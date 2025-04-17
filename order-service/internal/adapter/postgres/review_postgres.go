package postgres

import (
	"order-service/internal/model"
)

type ReviewRepository interface {
	Create(review *model.Review) error
	Update(review *model.Review) error
	Delete(id string) error
}

func (r *PostgresRepo) Create(review *model.Review) error {
	query := `INSERT INTO reviews (id, product_id, user_id, rating, comment, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query,
		review.ID,
		review.ProductID,
		review.UserID,
		review.Rating,
		review.Comment,
		review.CreatedAt,
		review.UpdatedAt,
	)
	return err
}

func (r *PostgresRepo) Update(review *model.Review) error {
	query := `UPDATE reviews SET rating = $1, comment = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, review.Rating, review.Comment, review.UpdatedAt, review.ID)
	return err
}
