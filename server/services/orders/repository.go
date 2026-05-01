package orders

import (
	"database/sql"

	"github.com/gitTanzj/server/types"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetOrders() ([]*types.Order, error) {
	rows, err := r.db.Query("SELECT * FROM orders;")
	if err != nil {
		return nil, err
	}

	orders := make([]*types.Order, 0)
	for rows.Next() {
		o, err := scanRowIntoOrder(rows)
		if err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	return orders, nil
}

func scanRowIntoOrder(rows *sql.Rows) (*types.Order, error) {
	order := new(types.Order)

	err := rows.Scan(
		&order.Id,
		&order.CustomerId,
		&order.Total,
		&order.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return order, nil
}
