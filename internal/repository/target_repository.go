package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"monitor-bot/internal/models"
)

type TargetRepository struct {
	db *sqlx.DB
}

func NewTargetRepository(db *sqlx.DB) *TargetRepository {
	return &TargetRepository{db: db}
}

func (r *TargetRepository) Create(ctx context.Context, t *models.Target) error {
	query := `
		INSERT INTO targets (title, description, deadline, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		t.Title, t.Description, t.Deadline,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TargetRepository) GetByID(ctx context.Context, id int64) (*models.Target, error) {
	var target models.Target
	err := r.db.GetContext(ctx, &target, "SELECT * FROM targets WHERE id = $1", id)
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
	t.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE targets SET title=$1, description=$2, deadline=$3, updated_at=$4 WHERE id=$5`,
		t.Title, t.Description, t.Deadline, t.UpdatedAt, t.ID,
	)
	return err
}

func (r *TargetRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM targets WHERE id=$1", id)
	return err
}
