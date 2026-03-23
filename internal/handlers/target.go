package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"monitor-bot/internal/models"
	"monitor-bot/internal/repository"
)

func GetAllTargets(repo *repository.TargetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		targets, err := repo.GetAll(ctx)
		if err != nil {
			http.Error(w, "Failed to get targets: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(targets)
	}
}

func GetTargetByID(repo *repository.TargetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid target ID", http.StatusBadRequest)
			return
		}

		target, err := repo.GetByID(ctx, id)
		if err != nil {
			http.Error(w, "Target not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(target)
	}
}

func CreateTarget(repo *repository.TargetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var t models.Target
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		t.CreatedAt = time.Now()

		if err := repo.Create(ctx, &t); err != nil {
			http.Error(w, "Failed to create target: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	}
}

func UpdateTarget(repo *repository.TargetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid target ID", http.StatusBadRequest)
			return
		}

		var t models.Target
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		t.ID = id
		if err := repo.Update(ctx, &t); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				http.Error(w, "Target not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to update target: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
	}
}

func DeleteTarget(repo *repository.TargetRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid target ID", http.StatusBadRequest)
			return
		}

		if err := repo.Delete(ctx, id); err != nil {
			http.Error(w, "Failed to delete target: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
