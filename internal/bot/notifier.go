package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (b *Bot) Notify(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.API.Send(msg)
	if err != nil {
		if tErr, ok := err.(tgbotapi.Error); ok {
			if tErr.Code == 403 {
				log.Printf("User %d заблокировал бота или удалил чат. Деактивируем пользователя.", chatID)
				user, getUserErr := b.UserService.GetByChatID(chatID)
				if getUserErr != nil {
					log.Printf("Не удалось найти пользователя по chatID %d: %v", chatID, getUserErr)
					return err
				}
				if err := b.UserService.Deactivate(user.ID); err != nil {
					log.Printf("Ошибка при деактивации пользователя %d: %v", chatID, err)
				}
			}
		} else {
			log.Printf("Ошибка отправки уведомления пользователю %d: %v", chatID, err)
		}
		return err
	}
	return nil
}
