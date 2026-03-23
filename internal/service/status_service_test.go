package service

import (
	"context"
	"errors"
	"testing"

	"monitor-bot/internal/models"
)

type stubBot struct {
	err error
}

func (b stubBot) Notify(chatID int64, text string) error {
	return b.err
}

type stubUserService struct {
	users             map[int64]*models.User
	deactivatedUserID int64
}

func (s *stubUserService) CreateOrActivate(chatID int64) error {
	return nil
}

func (s *stubUserService) Deactivate(userID int64) error {
	s.deactivatedUserID = userID
	return nil
}

func (s *stubUserService) GetByChatID(chatID int64) (*models.User, error) {
	return nil, errors.New("not implemented")
}

func (s *stubUserService) GetByID(userID int64) (*models.User, error) {
	user, ok := s.users[userID]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}

func TestStatusServiceSend_DeactivatesUserByUserIDOnInactiveError(t *testing.T) {
	userService := &stubUserService{}
	service := &StatusService{
		Bot:         stubBot{err: errors.New("Forbidden: bot was blocked by the user")},
		UserService: userService,
	}

	sub := models.Subscription{
		ID:     10,
		UserID: 42,
	}

	userService.users = map[int64]*models.User{
		sub.UserID: {
			ID:             sub.UserID,
			TelegramChatID: 777000,
		},
	}

	service.send(context.Background(), sub, "test message")

	if userService.deactivatedUserID != sub.UserID {
		t.Fatalf("expected deactivation by userID %d, got %d", sub.UserID, userService.deactivatedUserID)
	}
}
