package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"monitor-bot/internal/models"
)

type TargetRepository struct {
	db *sqlx.DB
}

// Конструктор
func NewTargetRepository(db *sqlx.DB) *TargetRepository {
	return &TargetRepository{db: db}
}

func (r *TargetRepository) Create(ctx context.Context, t *models.Target) error {
	query := `
		INSERT INTO targets
			(name, url, method, expected_status, body_regex, interval_seconds, timeout_seconds, region_id, created_by, enabled, type)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, created_at
	`
	if t.Type == "" {
		t.Type = "http"
	}

	if t.RegionID == 0 {
		t.RegionID = 1
	}

	return r.db.QueryRowContext(ctx, query,
		t.Name,
		t.URL,
		t.Method,
		t.ExpectedStatus,
		t.BodyRegex,
		t.IntervalSeconds,
		t.TimeoutSeconds,
		t.RegionID,
		t.CreatedBy,
		t.Enabled,
		t.Type,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *TargetRepository) GetByID(ctx context.Context, id int64) (*models.Target, error) {
	var target models.Target
	err := r.db.GetContext(ctx, &target, "SELECT * FROM targets WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *TargetRepository) GetByURL(ctx context.Context, url string) (*models.Target, error) {
	var target models.Target
	query := `SELECT * FROM targets WHERE url = $1 LIMIT 1`
	err := r.db.GetContext(ctx, &target, query, url)
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *TargetRepository) GetAll(ctx context.Context) ([]models.Target, error) {
	var targets []models.Target
	err := r.db.SelectContext(ctx, &targets, "SELECT * FROM targets ORDER BY id")
	return targets, err
}

func (r *TargetRepository) Update(ctx context.Context, t *models.Target) error {
	if t.Type == "" {
		t.Type = "http"
	}
	if t.RegionID == 0 {
		t.RegionID = 1
	}
	_, err := r.db.ExecContext(ctx, `
		UPDATE targets SET
			name=$1,
			url=$2,
			method=$3,
			expected_status=$4,
			body_regex=$5,
			interval_seconds=$6,
			timeout_seconds=$7,
			region_id=$8,
			enabled=$9,
		    type=$10
		WHERE id=$11
	`,
		t.Name,
		t.URL,
		t.Method,
		t.ExpectedStatus,
		t.BodyRegex,
		t.IntervalSeconds,
		t.TimeoutSeconds,
		t.RegionID,
		t.Enabled,
		t.Type,
		t.ID,
	)
	return err
}

func (r *TargetRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM targets WHERE id=$1", id)
	return err
}
