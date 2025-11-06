package service

type Bot interface {
	Notify(chatID int64, text string) error
}
