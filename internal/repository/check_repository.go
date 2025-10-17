package repository

import (
	"context"
	"monitor-bot/internal/models"

	"github.com/jmoiron/sqlx"
)

type CheckRepository struct {
	db *sqlx.DB
}

func NewCheckRepository(db *sqlx.DB) *CheckRepository {
	return &CheckRepository{db: db}
}

func (r *CheckRepository) Save(ctx context.Context, c *models.Check) error {
	query := `
		INSERT INTO checks (target_id, timestamp, status, http_code, response_time_ms, error, region)
		VALUES ($1, NOW(), $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		c.TargetID,
		c.Status,
		c.HttpCode,
		c.ResponseTimeMs,
		c.Error,
		c.Region,
	)
	return err
}
