package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"monitor-bot/internal/bot"
	"monitor-bot/internal/db"
	"monitor-bot/internal/repository"
	"monitor-bot/internal/scheduler"
	"monitor-bot/internal/service"
	"monitor-bot/internal/worker"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	notifyInterval    = 1 * time.Minute
	schedulerInterval = 10 * time.Second
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Не удалось загрузить .env (файл может отсутствовать)")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbConn := db.ConnectDB()
	defer dbConn.Close()

	userRepo := repository.NewUserRepository(dbConn)
	subRepo := repository.NewSubscriptionRepository(dbConn)
	targetRepo := repository.NewTargetRepository(dbConn)
	checkRepo := repository.NewCheckRepository(dbConn)

	userService := service.NewUserService(userRepo)
	subscriptionService := service.NewSubscriptionService(subRepo, userService, targetRepo)

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не задан")
	}

	telegramBot, err := bot.NewBot(telegramToken, userService, subscriptionService)
	if err != nil {
		log.Fatal("Ошибка создания Telegram-бота:", err)
	}

	statusService := &service.StatusService{
		CheckRepo:      checkRepo,
		SubRepo:        subRepo,
		TargetRepo:     targetRepo,
		Bot:            telegramBot,
		UserService:    userService,
		NotifyInterval: notifyInterval,
	}

	w := worker.NewWorker(targetRepo, checkRepo, statusService)

	s := scheduler.NewScheduler(targetRepo, w, schedulerInterval)
	log.Printf("Worker started: scheduler interval=%s, notify interval=%s", schedulerInterval, notifyInterval)
	go s.Start(ctx)

	<-ctx.Done()
	log.Println("Worker shutting down...")
}
