package repository

import (
	"context"
	"monitor-bot/internal/models"

	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) GetByTarget(ctx context.Context, targetID int64) ([]models.Subscription, error) {
	var subs []models.Subscription
	query := `
		SELECT s.id, s.user_id, s.target_id, s.notify_down_only, s.min_retries, s.last_notified, u.telegram_chat_id AS chat_id
		FROM subscriptions s
		JOIN users u ON u.id = s.user_id
		WHERE s.target_id = $1 AND u.is_active = TRUE
	`
	err := r.db.SelectContext(ctx, &subs, query, targetID)
	return subs, err
}

func (r *SubscriptionRepository) UpdateLastNotified(ctx context.Context, sub models.Subscription) error {
	query := `
		UPDATE subscriptions
		SET last_notified = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, sub.ID)
	return err
}

func (r *SubscriptionRepository) Subscribe(ctx context.Context, userID, targetID int64) error {
	query := `
		INSERT INTO subscriptions (user_id, target_id, notify_down_only, min_retries, created_at)
		VALUES ($1, $2, TRUE, 1, NOW())
		ON CONFLICT (user_id, target_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID, targetID)
	return err
}

func (r *SubscriptionRepository) Unsubscribe(ctx context.Context, userID, targetID int64) error {
	query := `
		DELETE FROM subscriptions
		WHERE user_id = $1 AND target_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, userID, targetID)
	return err
}
