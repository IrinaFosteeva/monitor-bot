package main

import (
	"log"
	"net/http"

	"monitor-bot/internal/db"
	"monitor-bot/internal/repository"
	"monitor-bot/internal/routes"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Не удалось загрузить .env (файл может отсутствовать)")
	}

	database := db.ConnectDB()
	defer database.Close()

	repo := repository.NewTargetRepository(database)
	r := routes.SetupRoutes(repo)

	log.Println("🚀 Server started on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("❌ Server error:", err)
	}
}
