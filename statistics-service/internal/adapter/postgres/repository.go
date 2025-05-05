package postgres

import (
	"database/sql"
	"log"
	"statistics-service/internal/model"
)

type StatisticsRepo struct {
	db *sql.DB
}

func NewStatisticsRepository(dsn string) (*StatisticsRepo, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &StatisticsRepo{db: db}, db.Ping()
}

func (r *StatisticsRepo) SaveOrderEvent(event *model.OrderEvent) error {
	_, err := r.db.Exec(`
		INSERT INTO order_events (order_id, user_id, total, timestamp)
		VALUES ($1, $2, $3, $4)`,
		event.OrderID, event.UserID, event.Total, event.Timestamp)
	if err != nil {
		log.Println("SaveOrderEvent error:", err)
	}
	return err
}

func (r *StatisticsRepo) SaveProductEvent(event *model.ProductEvent) error {
	_, err := r.db.Exec(`
		INSERT INTO product_events (product_id, name, category, price, stock, action, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		event.ProductID, event.Name, event.Category, event.Price, event.Stock, event.Action, event.Timestamp)
	if err != nil {
		log.Println("SaveProductEvent error:", err)
	}
	return err
}

// --- Методы для статистики ---

func (r *StatisticsRepo) CountOrdersByUser(userID string) (int, error) {
	row := r.db.QueryRow(`SELECT COUNT(*) FROM order_events WHERE user_id = $1`, userID)
	var count int
	return count, row.Scan(&count)
}

func (r *StatisticsRepo) OrdersGroupedByHour(userID string) (map[int]int, error) {
	rows, err := r.db.Query(`
		SELECT EXTRACT(HOUR FROM timestamp)::int AS hour, COUNT(*) 
		FROM order_events 
		WHERE user_id = $1 
		GROUP BY hour 
		ORDER BY hour`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]int)
	for rows.Next() {
		var hour, count int
		if err := rows.Scan(&hour, &count); err != nil {
			return nil, err
		}
		result[hour] = count
	}
	return result, nil
}

func (r *StatisticsRepo) CountTotalUsers() (int, error) {
	row := r.db.QueryRow(`SELECT COUNT(DISTINCT user_id) FROM order_events`)
	var count int
	return count, row.Scan(&count)
}

func (r *StatisticsRepo) CountTotalProducts() (int, error) {
	row := r.db.QueryRow(`SELECT COUNT(DISTINCT product_id) FROM product_events WHERE action = 'created'`)
	var count int
	return count, row.Scan(&count)
}
