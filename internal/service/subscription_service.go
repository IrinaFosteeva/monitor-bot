package service

import (
	"context"
	"errors"
	"monitor-bot/internal/models"
	"monitor-bot/internal/repository"
)

type SubscriptionServiceInterface interface {
	SubscribeByURL(chatID int64, url string) error
	UnsubscribeByURL(chatID int64, url string) error
}

type SubscriptionService struct {
	SubRepo    *repository.SubscriptionRepository
	UserRepo   UserServiceInterface
	TargetRepo repository.TargetRepositoryInterface
}

func NewSubscriptionService(subRepo *repository.SubscriptionRepository, userRepo UserServiceInterface, targetRepo repository.TargetRepositoryInterface) *SubscriptionService {
	return &SubscriptionService{
		SubRepo:    subRepo,
		UserRepo:   userRepo,
		TargetRepo: targetRepo,
	}
}

func (s *SubscriptionService) SubscribeByURL(chatID int64, url string) error {
	ctx := context.Background()

	user, err := s.UserRepo.GetByChatID(chatID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	target, err := s.TargetRepo.GetByURL(ctx, url)
	if err != nil {
		target = &models.Target{
			Name:            url,
			URL:             url,
			Method:          "GET",
			ExpectedStatus:  200,
			IntervalSeconds: 60,
			TimeoutSeconds:  5,
			RegionID:        1,
			Type:            "http",
			Enabled:         true,
		}
		if err := s.TargetRepo.Create(ctx, target); err != nil {
			return errors.New("не удалось создать цель")
		}
	}

	return s.SubRepo.Subscribe(ctx, user.ID, target.ID)
}

func (s *SubscriptionService) UnsubscribeByURL(chatID int64, url string) error {
	ctx := context.Background()

	user, err := s.UserRepo.GetByChatID(chatID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	target, err := s.TargetRepo.GetByURL(ctx, url)
	if err != nil {
		return errors.New("цель с таким URL не найдена")
	}

	return s.SubRepo.Unsubscribe(ctx, user.ID, target.ID)
}
