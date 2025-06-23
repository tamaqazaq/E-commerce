package postgres

import (
	"database/sql"
	"errors"
	"order-service/internal/model"
	"order-service/internal/usecase"
	"time"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepository(connStr string) (usecase.OrderRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresRepo{db: db}, nil
}

func (r *PostgresRepo) Save(order *model.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	orderQuery := `INSERT INTO orders (id, user_id, total, status, timestamp) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(orderQuery, order.ID, order.UserID, order.Total, order.Status, order.Timestamp)
	if err != nil {
		return err
	}

	itemQuery := `INSERT INTO order_items (id, order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4, $5)`
	for _, item := range order.Items {
		_, err := tx.Exec(itemQuery, item.ID, item.OrderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *PostgresRepo) FindByID(id string) (*model.Order, error) {
	orderQuery := `SELECT id, user_id, total, status, timestamp FROM orders WHERE id = $1`
	row := r.db.QueryRow(orderQuery, id)
	order := &model.Order{}
	if err := row.Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.Timestamp); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	itemQuery := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`
	rows, err := r.db.Query(itemQuery, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := model.OrderItem{}
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	return order, nil
}
func (r *PostgresRepo) UpdateStatus(id, status string) error {
	query := `UPDATE orders SET status = $1, timestamp = $2 WHERE id = $3`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

func (r *PostgresRepo) FindByUserID(userID string) ([]*model.Order, error) {
	ordersQuery := `SELECT id, user_id, total, status, timestamp FROM orders WHERE user_id = $1`
	rows, err := r.db.Query(ordersQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		order := &model.Order{}
		if err := rows.Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.Timestamp); err != nil {
			return nil, err
		}

		itemQuery := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`
		itemRows, err := r.db.Query(itemQuery, order.ID)
		if err != nil {
			return nil, err
		}
		for itemRows.Next() {
			item := model.OrderItem{}
			if err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}
		itemRows.Close()
		orders = append(orders, order)
	}
	return orders, nil
}
func (r *PostgresRepo) FindAll() ([]*model.Order, error) {
	query := `SELECT id, user_id, total, status, timestamp FROM orders`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		order := &model.Order{}
		if err := rows.Scan(&order.ID, &order.UserID, &order.Total, &order.Status, &order.Timestamp); err != nil {
			return nil, err
		}

		itemQuery := `SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`
		itemRows, err := r.db.Query(itemQuery, order.ID)
		if err != nil {
			return nil, err
		}
		for itemRows.Next() {
			item := model.OrderItem{}
			if err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
				return nil, err
			}
			order.Items = append(order.Items, item)
		}
		itemRows.Close()
		orders = append(orders, order)
	}
	return orders, nil
}
