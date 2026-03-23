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

func (s *UserService) Deactivate(userID int64) error {
	return s.repo.Deactivate(userID)
}

func (s *UserService) GetByChatID(chatID int64) (*models.User, error) {
	return s.repo.GetByChatID(chatID)
}

func (s *UserService) GetByID(userID int64) (*models.User, error) {
	return s.repo.GetByID(userID)
}
