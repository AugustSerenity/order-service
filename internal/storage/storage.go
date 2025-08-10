package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/AugustSerenity/order-service/internal/model"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveOrder(order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO orders (uid, content)
		VALUES ($1, $2)
		ON CONFLICT (uid) DO NOTHING
	`, order.OrderUID, data)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert order: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Storage) GetByID(id string) (model.Order, error) {
	var raw json.RawMessage

	err := s.db.QueryRow(`SELECT content FROM orders WHERE uid = $1`, id).Scan(&raw)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Order{}, fmt.Errorf("order not found: %w", err)
		}
		return model.Order{}, fmt.Errorf("query failed: %w", err)
	}

	var order model.Order
	if err := json.Unmarshal(raw, &order); err != nil {
		return model.Order{}, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return order, nil
}

func (s *Storage) GetAllOrders() ([]model.Order, error) {
	rows, err := s.db.Query(`SELECT content FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var raw json.RawMessage
		if err := rows.Scan(&raw); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var order model.Order
		if err := json.Unmarshal(raw, &order); err != nil {
			continue
		}

		orders = append(orders, order)
	}

	return orders, nil
}
