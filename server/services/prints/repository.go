package prints

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

func (r *Repository) GetPrints() ([]*types.Print, error) {
	rows, err := r.db.Query(`
		SELECT p.id, f.filament_name, p.path_to_file, p.filament_amount_used, p.estimated_print_time
		FROM prints p
		JOIN filaments f ON p.filament_id = f.id
	`)
	if err != nil {
		return nil, err
	}

	prints := make([]*types.Print, 0)
	for rows.Next() {
		p, err := scanRowIntoPrint(rows)
		if err != nil {
			return nil, err
		}

		prints = append(prints, p)
	}

	return prints, nil
}

func scanRowIntoPrint(rows *sql.Rows) (*types.Print, error) {
	print := new(types.Print)

	err := rows.Scan(
		&print.Id,
		&print.FilamentName,
		&print.PathToFile,
		&print.FilamentAmountUsed,
		&print.EstimatedPrintTime,
	)

	if err != nil {
		return nil, err
	}

	return print, nil
}
