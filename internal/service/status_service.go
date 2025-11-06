package service

import (
	"monitor-bot/internal/repository"
	"time"
)

type StatusService struct {
	CheckRepo      *repository.CheckRepository
	SubRepo        *repository.SubscriptionRepository
	TargetRepo     repository.TargetRepositoryInterface
	Bot            Bot
	NotifyInterval time.Duration
}

//func (s *StatusService) ProcessCheck(ctx context.Context, result *models.Check) error {
//	// Получаем последнюю проверку
//	prev, err := s.CheckRepo.GetLastByTarget(ctx, result.TargetID)
//	if err != nil && err != repository.ErrNotFound {
//		return fmt.Errorf("не удалось получить последнюю проверку: %w", err)
//	}
//
//	// Сохраняем текущую проверку
//	if err := s.CheckRepo.Save(ctx, result); err != nil {
//		return fmt.Errorf("не удалось сохранить проверку: %w", err)
//	}
//
//	// Получаем цель
//	target, err := s.TargetRepo.GetByID(ctx, result.TargetID)
//	if err != nil {
//		return fmt.Errorf("не удалось получить target: %w", err)
//	}
//	if target == nil {
//		return fmt.Errorf("target с id=%d не найден", result.TargetID)
//	}
//
//	// Получаем подписки
//	subs, err := s.SubRepo.GetByTarget(ctx, result.TargetID)
//	if err != nil {
//		return fmt.Errorf("не удалось получить подписки: %w", err)
//	}
//
//	for _, sub := range subs {
//		// Если подписчик хочет уведомления только о падениях
//		if result.Status == "down" {
//			downCount, err := s.CheckRepo.CountLastNDown(ctx, result.TargetID, sub.MinRetries)
//			if err != nil {
//				return fmt.Errorf("ошибка подсчета последних падений: %w", err)
//			}
//			if downCount < sub.MinRetries {
//				continue
//			}
//		}
//
//		// Проверяем интервал уведомлений
//		if !sub.LastNotified.IsZero() && time.Since(sub.LastNotified) < s.NotifyInterval {
//			continue
//		}
//
//		statusText := "UP ✅"
//		if result.Status == "down" {
//			statusText = "DOWN ❌"
//		}
//		msg := fmt.Sprintf("Цель %s изменила состояние: %s", target.Name, statusText)
//
//		// Отправляем уведомление
//		err := s.Bot.Notify(sub.ChatID, msg)
//		if err != nil {
//			// Если чат удален или бот заблокирован
//			if isUserInactiveError(err) {
//				if err := s.SubRepo.DeactivateUser(sub.UserID); err != nil {
//					fmt.Println("не удалось деактивировать пользователя:", err)
//				}
//				fmt.Printf("Пользователь %d заблокировал бота или удалил чат\n", sub.UserID)
//			} else {
//				fmt.Println("Ошибка отправки уведомления:", err)
//			}
//			continue
//		}
//
//		// Обновляем LastNotified
//		sub.LastNotified = time.Now()
//		if err := s.SubRepo.UpdateLastNotified(ctx, sub); err != nil {
//			fmt.Println("не удалось обновить LastNotified:", err)
//		}
//	}
//
//	return nil
//}
//
//// Вспомогательная функция для определения ошибок Telegram
//func isUserInactiveError(err error) bool {
//	if tErr, ok := err.(interface{ Error() string }); ok {
//		msg := tErr.Error()
//		return msg == "Forbidden: bot was blocked by the user" || msg == "Forbidden: chat not found"
//	}
//	return false
//}
