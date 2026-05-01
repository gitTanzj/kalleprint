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

func (r *Repository) GetFilaments() ([]*types.Filament, error) {
	rows, err := r.db.Query("SELECT BIN_TO_UUID(id), filament_name, current_amount, cost_per_gram FROM filaments")
	if err != nil {
		return nil, err
	}

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
	rows, err := r.db.Query("SELECT BIN_TO_UUID(id), filament_name, current_amount, cost_per_gram FROM filaments WHERE id = UUID_TO_BIN(?)", id)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("filament not found")
	}

	return scanRowIntoFilament(rows)
}

func scanRowIntoFilament(rows *sql.Rows) (*types.Filament, error) {
	filament := new(types.Filament)

	err := rows.Scan(
		&filament.Id,
		&filament.FilamentName,
		&filament.CurrentAmount,
		&filament.CostPerGram,
	)

	if err != nil {
		return nil, err
	}

	return filament, nil
}
