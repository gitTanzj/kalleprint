package jobs

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

func (r *Repository) GetJobs() ([]*types.Job, error) {
	rows, err := r.db.Query(`
		SELECT BIN_TO_UUID(j.id), BIN_TO_UUID(j.order_id), BIN_TO_UUID(j.filament_id),
		       COALESCE(f.filament_name, ''), j.gcode_path, j.status, j.created_at
		FROM jobs j
		LEFT JOIN filaments f ON j.filament_id = f.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]*types.Job, 0)
	for rows.Next() {
		j, err := scanRowIntoJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (r *Repository) GetJobByID(id string) (*types.Job, error) {
	row := r.db.QueryRow(`
		SELECT BIN_TO_UUID(j.id), BIN_TO_UUID(j.order_id), BIN_TO_UUID(j.filament_id),
		       COALESCE(f.filament_name, ''), j.gcode_path, j.status, j.created_at
		FROM jobs j
		LEFT JOIN filaments f ON j.filament_id = f.id
		WHERE j.id = UUID_TO_BIN(?)`, id)

	j := new(types.Job)
	err := row.Scan(&j.Id, &j.OrderId, &j.FilamentId, &j.FilamentName, &j.GcodePath, &j.Status, &j.CreatedAt)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (r *Repository) CreateJob(orderID, filamentID, gcodePath string) (*types.Job, error) {
	var newID string
	if err := r.db.QueryRow("SELECT UUID()").Scan(&newID); err != nil {
		return nil, err
	}

	_, err := r.db.Exec(
		`INSERT INTO jobs (id, order_id, filament_id, gcode_path)
		 VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), UUID_TO_BIN(?), ?)`,
		newID, orderID, filamentID, gcodePath,
	)
	if err != nil {
		return nil, err
	}

	return r.GetJobByID(newID)
}

func (r *Repository) UpdateJobStatus(id, status string) error {
	_, err := r.db.Exec(
		`UPDATE jobs SET status = ? WHERE id = UUID_TO_BIN(?)`,
		status, id,
	)
	return err
}

func scanRowIntoJob(rows *sql.Rows) (*types.Job, error) {
	j := new(types.Job)
	err := rows.Scan(&j.Id, &j.OrderId, &j.FilamentId, &j.FilamentName, &j.GcodePath, &j.Status, &j.CreatedAt)
	if err != nil {
		return nil, err
	}
	return j, nil
}
