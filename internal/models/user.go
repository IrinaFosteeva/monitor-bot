package models

import "time"

type User struct {
	ID             int64     `db:"id"`
	TelegramChatID int64     `db:"telegram_chat_id"`
	IsActive       bool      `db:"is_active"`
	IsPremium      bool      `db:"is_premium"`
	CreatedAt      time.Time `db:"created_at"`
}
