package service

import "monitor-bot/internal/models"

type UserServiceInterface interface {
	CreateOrActivate(chatID int64) error
	Deactivate(chatID int64) error
	GetByChatID(chatID int64) (*models.User, error)
}
