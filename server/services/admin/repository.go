package admin

import (
	"database/sql"
	"fmt"

	"github.com/gitTanzj/server/types"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAdminByEmail(email string) (*types.AdminUser, error) {
	row := r.db.QueryRow(
		"SELECT BIN_TO_UUID(id), email, password_hash, created_at FROM admin_users WHERE email = ?",
		email,
	)
	u := new(types.AdminUser)
	if err := row.Scan(&u.Id, &u.Email, &u.PasswordHash, &u.CreatedAt); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repository) GetOrders(status string) ([]*types.Order, error) {
	query := `SELECT BIN_TO_UUID(id), name, email, shipping_address,
	                 COALESCE(notes, ''), COALESCE(BIN_TO_UUID(print_id), ''), total, status, created_at
	          FROM orders`
	var args []any
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]*types.Order, 0)
	for rows.Next() {
		o := new(types.Order)
		if err := rows.Scan(
			&o.Id, &o.Name, &o.Email, &o.ShippingAddress,
			&o.Notes, &o.PrintId, &o.Total, &o.Status, &o.CreatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *Repository) GetOrderByID(id string) (*types.Order, error) {
	row := r.db.QueryRow(`
		SELECT BIN_TO_UUID(id), name, email, shipping_address,
		       COALESCE(notes, ''), COALESCE(BIN_TO_UUID(print_id), ''), total, status, created_at
		FROM orders WHERE id = UUID_TO_BIN(?)`, id)
	o := new(types.Order)
	if err := row.Scan(
		&o.Id, &o.Name, &o.Email, &o.ShippingAddress,
		&o.Notes, &o.PrintId, &o.Total, &o.Status, &o.CreatedAt,
	); err != nil {
		return nil, err
	}
	return o, nil
}

func (r *Repository) UpdateOrderStatus(id, status string) error {
	res, err := r.db.Exec(
		"UPDATE orders SET status = ? WHERE id = UUID_TO_BIN(?)", status, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

func (r *Repository) GetFilaments() ([]*types.Filament, error) {
	rows, err := r.db.Query(
		"SELECT BIN_TO_UUID(id), filament_name, current_amount, cost_per_gram, color_hex, ini_file_path FROM filaments",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	filaments := make([]*types.Filament, 0)
	for rows.Next() {
		f := new(types.Filament)
		if err := rows.Scan(&f.Id, &f.FilamentName, &f.CurrentAmount, &f.CostPerGram, &f.ColorHex, &f.IniFilePath); err != nil {
			return nil, err
		}
		filaments = append(filaments, f)
	}
	return filaments, nil
}

func (r *Repository) CreateFilament(name string, amountGrams int, costPerGram float64) (*types.Filament, error) {
	var newID string
	if err := r.db.QueryRow("SELECT UUID()").Scan(&newID); err != nil {
		return nil, err
	}
	_, err := r.db.Exec(
		"INSERT INTO filaments (id, filament_name, current_amount, cost_per_gram) VALUES (UUID_TO_BIN(?), ?, ?, ?)",
		newID, name, amountGrams, costPerGram,
	)
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(
		"SELECT BIN_TO_UUID(id), filament_name, current_amount, cost_per_gram, color_hex, ini_file_path FROM filaments WHERE id = UUID_TO_BIN(?)",
		newID,
	)
	f := new(types.Filament)
	if err := row.Scan(&f.Id, &f.FilamentName, &f.CurrentAmount, &f.CostPerGram, &f.ColorHex, &f.IniFilePath); err != nil {
		return nil, err
	}
	return f, nil
}

func (r *Repository) UpdateFilament(id, name string, amountGrams int, costPerGram float64) (*types.Filament, error) {
	res, err := r.db.Exec(
		"UPDATE filaments SET filament_name = ?, current_amount = ?, cost_per_gram = ? WHERE id = UUID_TO_BIN(?)",
		name, amountGrams, costPerGram, id,
	)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, fmt.Errorf("filament not found")
	}
	row := r.db.QueryRow(
		"SELECT BIN_TO_UUID(id), filament_name, current_amount, cost_per_gram, color_hex, ini_file_path FROM filaments WHERE id = UUID_TO_BIN(?)",
		id,
	)
	f := new(types.Filament)
	if err := row.Scan(&f.Id, &f.FilamentName, &f.CurrentAmount, &f.CostPerGram, &f.ColorHex, &f.IniFilePath); err != nil {
		return nil, err
	}
	return f, nil
}

func (r *Repository) DeleteFilament(id string) error {
	res, err := r.db.Exec("DELETE FROM filaments WHERE id = UUID_TO_BIN(?)", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("filament not found")
	}
	return nil
}

func (r *Repository) GetDashboardStats() (*types.DashboardStats, error) {
	row := r.db.QueryRow(`
		SELECT
			COUNT(*) AS total_orders,
			SUM(status = 'pending') AS pending,
			SUM(status IN ('printing', 'waiting-for-shipment')) AS in_progress,
			COALESCE(SUM(total), 0) AS total_revenue
		FROM orders
	`)
	s := new(types.DashboardStats)
	if err := row.Scan(&s.TotalOrders, &s.Pending, &s.InProgress, &s.TotalRevenue); err != nil {
		return nil, err
	}
	return s, nil
}
