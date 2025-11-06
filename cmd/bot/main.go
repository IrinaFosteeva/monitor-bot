package main

import (
	"github.com/joho/godotenv"
	"log"
	"monitor-bot/internal/bot"
	"monitor-bot/internal/db"
	"monitor-bot/internal/repository"
	"monitor-bot/internal/service"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Не удалось загрузить .env (файл может отсутствовать)")
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	database := db.ConnectDB()
	defer database.Close()

	userRepo := repository.NewUserRepository(database)
	subRepo := repository.NewSubscriptionRepository(database)
	targetRepo := repository.NewTargetRepository(database)

	userService := service.NewUserService(userRepo)
	subscriptionService := service.NewSubscriptionService(subRepo, userService, targetRepo)

	b, err := bot.NewBot(token, userService, subscriptionService)
	if err != nil {
		log.Fatal("Ошибка создания бота:", err)
	}

	log.Println("Bot started...")
	b.Start()
}
