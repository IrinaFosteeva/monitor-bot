package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"time"

	"monitor-bot/internal/db"
	"monitor-bot/internal/repository"
	"monitor-bot/internal/scheduler"
	"monitor-bot/internal/worker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Не удалось загрузить .env (файл может отсутствовать)")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn := db.ConnectDB()

	targetRepo := repository.NewTargetRepository(dbConn)
	checkRepo := repository.NewCheckRepository(dbConn)
	w := worker.NewWorker(targetRepo, checkRepo)

	s := scheduler.NewScheduler(targetRepo, w, 10*time.Second)
	go s.Start(ctx)

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Shutting down...")
}
