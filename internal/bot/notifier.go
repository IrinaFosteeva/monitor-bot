package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (b *Bot) NotifyUser(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := b.API.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки уведомления:", err)
		if err.Error() == "Forbidden: bot was blocked by the user" || err.Error() == "Forbidden: chat not found" {

			b.UserService.Deactivate(chatID)
		}
	}
}
