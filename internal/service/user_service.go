package service

import (
	"monitor-bot/internal/models"
	"monitor-bot/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateOrActivate(chatID int64) error {
	return s.repo.CreateOrActivate(chatID)
}

func (s *UserService) Deactivate(chatID int64) error {
	return s.repo.Deactivate(chatID)
}

func (s *UserService) GetByChatID(chatID int64) (*models.User, error) {
	return s.repo.GetByChatID(chatID)
}
