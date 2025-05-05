package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"inventory-service/internal/model"
)

type DiscountRepository interface {
	Save(discount *model.Discount) error
	GetProductsWithDiscount() ([]*model.Product, error)
	Delete(id string) error
}

type PostgresDiscountRepo struct {
	db *sql.DB
}

func NewPostgresDiscountRepository(db *sql.DB) *PostgresDiscountRepo {
	return &PostgresDiscountRepo{db: db}
}

func (r *PostgresDiscountRepo) Save(discount *model.Discount) error {
	productsJSON, err := json.Marshal(discount.ApplicableProducts)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO discounts (id, name, description, discount_percentage, applicable_products, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = r.db.Exec(query,
		discount.ID,
		discount.Name,
		discount.Description,
		discount.DiscountPercentage,
		productsJSON,
		discount.StartDate,
		discount.EndDate,
		discount.IsActive,
	)
	return err
}

func (r *PostgresDiscountRepo) GetProductsWithDiscount() ([]*model.Product, error) {
	query := `
		SELECT p.id, p.name, p.category, p.price, p.stock
		FROM products p
		JOIN discounts d ON d.is_active = TRUE
		WHERE p.id = ANY (
			SELECT jsonb_array_elements_text(applicable_products)::text FROM discounts WHERE is_active = TRUE
		)
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		product := &model.Product{}
		if err := rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *PostgresDiscountRepo) Delete(id string) error {
	query := `DELETE FROM discounts WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("discount not found")
	}
	return nil
}
