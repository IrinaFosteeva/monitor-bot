package repository

import (
	"context"
	"database/sql"
	"errors"
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
		INSERT INTO checks (target_id, timestamp, status, http_code, response_time_ms, error)
		VALUES ($1, NOW(), $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		c.TargetID,
		c.Status,
		c.HttpCode,
		c.ResponseTimeMs,
		c.Error,
	)
	return err
}

func (r *CheckRepository) GetLastByTarget(ctx context.Context, targetID int64) (*models.Check, error) {
	var c models.Check
	query := `
		SELECT *
		FROM checks
		WHERE target_id = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &c, query, targetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *CheckRepository) CountLastNDown(ctx context.Context, targetID int64, n int) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM (
			SELECT status
			FROM checks
			WHERE target_id = $1
			ORDER BY timestamp DESC
			LIMIT $2
		) AS recent_checks
		WHERE status = 'down'
	`
	err := r.db.GetContext(ctx, &count, query, targetID, n)
	return count, err
}
