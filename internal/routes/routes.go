package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"monitor-bot/internal/handlers"
	"monitor-bot/internal/repository"
)

func SetupRoutes(repo *repository.TargetRepository) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/targets", func(r chi.Router) {
		r.Get("/", handlers.GetAllTargets(repo))
		r.Post("/", handlers.CreateTarget(repo))
		r.Get("/{id}", handlers.GetTargetByID(repo))
		r.Put("/{id}", handlers.UpdateTarget(repo))
		r.Delete("/{id}", handlers.DeleteTarget(repo))
	})

	return r
}
