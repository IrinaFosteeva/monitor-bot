package service

import "monitor-bot/internal/models"

type UserServiceInterface interface {
	CreateOrActivate(chatID int64) error
	Deactivate(userID int64) error
	GetByID(userID int64) (*models.User, error)
	GetByChatID(chatID int64) (*models.User, error)
}
