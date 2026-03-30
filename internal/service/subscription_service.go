package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"monitor-bot/internal/models"
	"monitor-bot/internal/repository"
)

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

func (s *SubscriptionService) SubscribeByURL(ctx context.Context, chatID int64, url string) error {
	if !IsValidURL(url) {
		return errors.New("❌ Невалидный URL. Используйте формат: https://example.com")
	}

	user, err := s.UserRepo.GetByChatID(chatID)
	if err != nil {
		return errors.New("❌ Пользователь не найден. Используйте /start для регистрации")
	}

	subsCount, err := s.SubRepo.CountByUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("❌ Ошибка проверки лимита: %w", err)
	}

	maxSubs := GetMaxSubscriptions(user)
	if subsCount >= maxSubs {
		if user.IsPremium {
			return fmt.Errorf("❌ Лимит подписок: %d/%d", subsCount, maxSubs)
		}
		return fmt.Errorf("❌ Лимит подписок: %d/%d\n\n💎 Upgrade до Premium для 10 целей!", subsCount, maxSubs)
	}

	target, err := s.TargetRepo.GetByURL(ctx, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			target = &models.Target{
				Name:            url,
				URL:             url,
				Method:          "GET",
				ExpectedStatus:  200,
				IntervalSeconds: 60,
				TimeoutSeconds:  10,
				RegionID:        1,
				Type:            "http",
				Enabled:         true,
			}
			if err := s.TargetRepo.Create(ctx, target); err != nil {
				return fmt.Errorf("❌ Не удалось создать цель: %w", err)
			}
		} else {
			return fmt.Errorf("❌ Ошибка поиска цели: %w", err)
		}
	} else {
		// Цель найдена, если она отключена - включаем обратно
		if !target.Enabled {
			if err := s.TargetRepo.Enable(ctx, target.ID); err != nil {
				return fmt.Errorf("❌ Не удалось активировать цель: %w", err)
			}
		}
	}

	if err := s.SubRepo.Subscribe(ctx, user.ID, target.ID); err != nil {
		return fmt.Errorf("❌ Не удалось создать подписку: %w", err)
	}

	return nil
}

func (s *SubscriptionService) UnsubscribeByURL(ctx context.Context, chatID int64, url string) error {
	user, err := s.UserRepo.GetByChatID(chatID)
	if err != nil {
		return errors.New("❌ Пользователь не найден")
	}

	target, err := s.TargetRepo.GetByURL(ctx, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("❌ Цель с таким URL не найдена")
		}
		return fmt.Errorf("❌ Ошибка поиска цели: %w", err)
	}

	if err := s.SubRepo.Unsubscribe(ctx, user.ID, target.ID); err != nil {
		return fmt.Errorf("❌ Не удалось отписаться: %w", err)
	}

	subs, err := s.SubRepo.GetByTarget(ctx, target.ID)
	if err != nil {
		return fmt.Errorf("⚠️ Отписка выполнена, но ошибка проверки подписчиков: %w", err)
	}

	if len(subs) == 0 {
		if err := s.TargetRepo.Disable(ctx, target.ID); err != nil {
			return fmt.Errorf("⚠️ Отписка выполнена, но не удалось отключить цель: %w", err)
		}
	}

	return nil
}
