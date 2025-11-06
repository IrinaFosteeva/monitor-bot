package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"monitor-bot/internal/service"
)

type Bot struct {
	API                 *tgbotapi.BotAPI
	UserService         service.UserServiceInterface
	SubscriptionService service.SubscriptionServiceInterface
}

func NewBot(
	token string,
	userService service.UserServiceInterface,
	subService service.SubscriptionServiceInterface,
) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{
		API:                 botAPI,
		UserService:         userService,
		SubscriptionService: subService,
	}, nil
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Command() {
			case "start":
				b.handleStart(update.Message.Chat.ID)
			case "test":
				b.Notify(update.Message.Chat.ID, "Тестовое уведомление: цель недоступна!")
			case "subscribe":
				url := update.Message.CommandArguments()
				if url == "" {
					b.Notify(update.Message.Chat.ID, "Пожалуйста, укажите URL после команды /subscribe")
					continue
				}
				err := b.SubscriptionService.SubscribeByURL(update.Message.Chat.ID, url)
				if err != nil {
					b.Notify(update.Message.Chat.ID, "Ошибка подписки: "+err.Error())
				} else {
					b.Notify(update.Message.Chat.ID, "Вы подписаны на "+url)
				}
			case "unsubscribe":
				url := update.Message.CommandArguments()
				if url == "" {
					b.Notify(update.Message.Chat.ID, "Пожалуйста, укажите URL после команды /unsubscribe")
					continue
				}
				err := b.SubscriptionService.UnsubscribeByURL(update.Message.Chat.ID, url)
				if err != nil {
					b.Notify(update.Message.Chat.ID, "Ошибка отписки: "+err.Error())
				} else {
					b.Notify(update.Message.Chat.ID, "Вы отписаны от "+url)
				}
			default:
				b.API.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда не распознана!"))
			}
		}
	}
}

func (b *Bot) handleStart(chatID int64) {
	err := b.UserService.CreateOrActivate(chatID)
	if err != nil {
		log.Println("Ошибка регистрации пользователя:", err)
		_, _ = b.API.Send(tgbotapi.NewMessage(chatID, "Ошибка при регистрации"))
		return
	}
	_, _ = b.API.Send(tgbotapi.NewMessage(chatID, "Регистрация прошла успешно!"))
}
