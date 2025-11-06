package models

import "time"

type Subscription struct {
	ID             int64      `db:"id" json:"id"`
	UserID         int64      `db:"user_id" json:"user_id"`
	TargetID       int64      `db:"target_id" json:"target_id"`
	NotifyDownOnly bool       `db:"notify_down_only" json:"notify_down_only"`
	MinRetries     int        `db:"min_retries" json:"min_retries"`
	LastNotified   *time.Time `db:"last_notified" json:"last_notified"`
	ChatID         int64      `db:"chat_id" json:"chat_id"`
}
