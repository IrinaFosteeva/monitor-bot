package repository

import (
	"github.com/jmoiron/sqlx"
	"monitor-bot/internal/models"
	"time"
)

type UserRepository struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateOrActivate(chatID int64) error {
	_, err := r.DB.Exec(`
		INSERT INTO users (telegram_chat_id, is_active, created_at)
		VALUES ($1, true, $2)
		ON CONFLICT (telegram_chat_id) DO UPDATE SET is_active = true
	`, chatID, time.Now())
	return err
}

func (r *UserRepository) Deactivate(chatID int64) error {
	_, err := r.DB.Exec(
		`UPDATE users SET is_active = false WHERE telegram_chat_id = $1`,
		chatID,
	)
	return err
}

func (r *UserRepository) GetByChatID(chatID int64) (*models.User, error) {
	var user models.User
	err := r.DB.Get(&user, `
		SELECT id, telegram_chat_id, is_active, created_at
		FROM users
		WHERE telegram_chat_id = $1
	`, chatID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
