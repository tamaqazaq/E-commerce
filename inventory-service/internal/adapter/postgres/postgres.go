package postgres

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"inventory-service/internal/model"
)

type PostgresProductRepo struct {
	db *sql.DB
}

func NewPostgresProductRepository(connStr string) (*PostgresProductRepo, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresProductRepo{db: db}, nil
}

func (r *PostgresProductRepo) Save(product *model.Product) error {
	query := `INSERT INTO products (id, name, category, price, stock) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, product.ID, product.Name, product.Category, product.Price, product.Stock)
	return err
}

func (r *PostgresProductRepo) FindByID(id string) (*model.Product, error) {
	query := `SELECT id, name, category, price, stock FROM products WHERE id = $1`
	row := r.db.QueryRow(query, id)
	product := &model.Product{}
	if err := row.Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.Stock); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (r *PostgresProductRepo) Update(product *model.Product) error {
	query := `UPDATE products SET name = $1, category = $2, price = $3, stock = $4 WHERE id = $5`
	_, err := r.db.Exec(query, product.Name, product.Category, product.Price, product.Stock, product.ID)
	return err
}

func (r *PostgresProductRepo) Delete(id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PostgresProductRepo) FindAll() ([]*model.Product, error) {
	query := `SELECT id, name, category, price, stock FROM products`
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

func (r *PostgresProductRepo) GetDB() *sql.DB {
	return r.db
}
