package filaments

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
