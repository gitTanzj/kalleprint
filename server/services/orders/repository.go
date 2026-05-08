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
	rows, err := r.db.Query(`
		SELECT BIN_TO_UUID(id), name, email, shipping_address,
		       COALESCE(notes, ''), COALESCE(BIN_TO_UUID(print_id), ''), total, status, created_at
		FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r *Repository) CreateAndReturnOrder(
	name string,
	email string,
	shippingAddress string,
	notes string,
	printID string,
	total float64,
) (*types.Order, error) {
	var newID string
	if err := r.db.QueryRow("SELECT UUID()").Scan(&newID); err != nil {
		return nil, err
	}

	_, err := r.db.Exec(
		`INSERT INTO orders (id, name, email, notes, shipping_address, print_id, total)
		VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, IF(? = '', NULL, UUID_TO_BIN(?)), ?)`,
		newID, name, email, notes, shippingAddress, printID, printID, total,
	)
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(`
		SELECT BIN_TO_UUID(id), name, email, shipping_address,
		       COALESCE(notes, ''), COALESCE(BIN_TO_UUID(print_id), ''), total, status, created_at
		FROM orders WHERE id = UUID_TO_BIN(?)`, newID)

	order := new(types.Order)
	if err := row.Scan(
		&order.Id, &order.Name, &order.Email, &order.ShippingAddress,
		&order.Notes, &order.PrintId, &order.Total, &order.Status, &order.CreatedAt,
	); err != nil {
		return nil, err
	}
	return order, nil
}

func scanRowIntoOrder(rows *sql.Rows) (*types.Order, error) {
	order := new(types.Order)
	err := rows.Scan(
		&order.Id, &order.Name, &order.Email, &order.ShippingAddress,
		&order.Notes, &order.PrintId, &order.Total, &order.Status, &order.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return order, nil
}
