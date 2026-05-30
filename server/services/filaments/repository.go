package filaments

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gitTanzj/server/types"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const filamentColumns = "BIN_TO_UUID(id), filament_name, current_amount, cost_per_gram, color_hex, ini_file_path"

func (r *Repository) GetFilaments() ([]*types.Filament, error) {
	rows, err := r.db.Query("SELECT " + filamentColumns + " FROM filaments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	filaments := make([]*types.Filament, 0)
	for rows.Next() {
		f, err := scanRowIntoFilament(rows)
		if err != nil {
			return nil, err
		}

		filaments = append(filaments, f)
	}

	return filaments, nil
}

func (r *Repository) GetFilamentByID(id string) (*types.Filament, error) {
	rows, err := r.db.Query("SELECT "+filamentColumns+" FROM filaments WHERE id = UUID_TO_BIN(?)", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("filament not found")
	}

	return scanRowIntoFilament(rows)
}

func (r *Repository) CreateFilament(id, name string, currentAmount int, costPerGram float64, colorHex, iniFilePath string) (*types.Filament, error) {
	if _, err := r.db.Exec(
		"INSERT INTO filaments (id, filament_name, current_amount, cost_per_gram, color_hex, ini_file_path) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?)",
		id, name, currentAmount, costPerGram, colorHex, iniFilePath,
	); err != nil {
		return nil, err
	}

	return r.GetFilamentByID(id)
}

func (r *Repository) UpdateFilament(id string, name *string, currentAmount *int, costPerGram *float64, colorHex *string, iniFilePath *string) error {
	var setClauses []string
	var args []any

	if name != nil {
		setClauses = append(setClauses, "filament_name = ?")
		args = append(args, *name)
	}
	if currentAmount != nil {
		setClauses = append(setClauses, "current_amount = ?")
		args = append(args, *currentAmount)
	}
	if costPerGram != nil {
		setClauses = append(setClauses, "cost_per_gram = ?")
		args = append(args, *costPerGram)
	}
	if colorHex != nil {
		setClauses = append(setClauses, "color_hex = ?")
		args = append(args, *colorHex)
	}
	if iniFilePath != nil {
		setClauses = append(setClauses, "ini_file_path = ?")
		args = append(args, *iniFilePath)
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, id)
	q := "UPDATE filaments SET " + strings.Join(setClauses, ", ") + " WHERE id = UUID_TO_BIN(?)"
	_, err := r.db.Exec(q, args...)
	return err
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

func scanRowIntoFilament(rows *sql.Rows) (*types.Filament, error) {
	filament := new(types.Filament)

	err := rows.Scan(
		&filament.Id,
		&filament.FilamentName,
		&filament.CurrentAmount,
		&filament.CostPerGram,
		&filament.ColorHex,
		&filament.IniFilePath,
	)

	if err != nil {
		return nil, err
	}

	return filament, nil
}
