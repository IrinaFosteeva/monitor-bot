package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"monitor-bot/internal/models"
	"monitor-bot/internal/repository"
	"time"
)

type StatusService struct {
	CheckRepo      *repository.CheckRepository
	SubRepo        *repository.SubscriptionRepository
	TargetRepo     repository.TargetRepositoryInterface
	Bot            Bot
	NotifyInterval time.Duration
	UserService    UserServiceInterface
}

func (s *StatusService) ProcessCheck(ctx context.Context, result *models.Check, target *models.Target) error {
	prev, err := s.CheckRepo.GetLastByTarget(ctx, result.TargetID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("ошибка получения предыдущей проверки: %w", err)
	}

	if err := s.CheckRepo.Save(ctx, result); err != nil {
		return fmt.Errorf("ошибка сохранения проверки: %w", err)
	}

	subs, err := s.SubRepo.GetByTarget(ctx, result.TargetID)
	if err != nil {
		return fmt.Errorf("ошибка получения подписок: %w", err)
	}

	for _, sub := range subs {
		s.handleNotification(ctx, sub, prev, result, target)
	}

	return nil
}

func (s *StatusService) handleNotification(ctx context.Context, sub models.Subscription, prev *models.Check, current *models.Check, target *models.Target) {
	prevStatus := ""
	if prev != nil {
		prevStatus = prev.Status
	}
	newStatus := current.Status
	if sub.NotifyDownOnly && newStatus == "up" {
		return
	}

	log.Printf("TARGET %s (%d): prevStatus=%q, newStatus=%q\n", target.Name, target.ID, prevStatus, newStatus)

	switch {
	case prevStatus == "up" && newStatus == "down":
		if !s.shouldNotifyDown(ctx, sub, current.TargetID) {
			return
		}
		msg := fmt.Sprintf("⚠️ Цель *%s* недоступна!\nURL: %s", target.Name, target.URL)
		s.send(ctx, sub, msg)

	case prevStatus == "down" && newStatus == "down":
		if sub.LastNotified != nil && time.Since(*sub.LastNotified) < s.NotifyInterval {
			return
		}
		msg := fmt.Sprintf("⏳ Цель *%s* всё ещё недоступна\nURL: %s", target.Name, target.URL)
		s.send(ctx, sub, msg)

	case prevStatus == "down" && newStatus == "up":
		msg := fmt.Sprintf("✅ Цель *%s* восстановлена!\nURL: %s", target.Name, target.URL)
		s.send(ctx, sub, msg)
	}
}

func (s *StatusService) shouldNotifyDown(ctx context.Context, sub models.Subscription, targetID int64) bool {
	downCount, err := s.CheckRepo.CountLastNDown(ctx, targetID, sub.MinRetries)
	if err != nil {
		log.Println("Ошибка подсчета порога падений:", err)
		return false
	}
	return downCount >= sub.MinRetries
}

func (s *StatusService) send(ctx context.Context, sub models.Subscription, message string) {
	err := s.Bot.Notify(sub.ChatID, message)
	if err != nil {
		if isUserInactiveError(err) {
			log.Printf("Пользователь %d заблокировал бота. Деактивируем.", sub.UserID)
			if deactErr := s.UserService.Deactivate(sub.UserID); deactErr != nil {
				log.Println("Ошибка деактивации пользователя:", deactErr)
			}
		} else {
			log.Println("Ошибка отправки уведомления:", err)
		}
		return
	}

	now := time.Now()
	sub.LastNotified = &now
	if err := s.SubRepo.UpdateLastNotified(ctx, sub); err != nil {
		log.Println("Ошибка обновления LastNotified:", err)
	}
}

func isUserInactiveError(err error) bool {
	if err == nil {
		return false
	}
	txt := err.Error()
	return txt == "Forbidden: bot was blocked by the user" || txt == "Forbidden: chat not found"
}
