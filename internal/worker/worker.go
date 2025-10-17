package worker

import (
	"context"
	"log"
	"monitor-bot/internal/models"
	"monitor-bot/internal/repository"
	"monitor-bot/internal/worker/checkers"
)

type Worker struct {
	targetRepo *repository.TargetRepository
	checkRepo  *repository.CheckRepository
}

func NewWorker(tRepo *repository.TargetRepository, cRepo *repository.CheckRepository) *Worker {
	return &Worker{targetRepo: tRepo, checkRepo: cRepo}
}

func (w *Worker) Run(ctx context.Context, target models.Target) {
	var (
		status   string
		code     *int
		duration int64
		errMsg   *string
	)

	switch target.Type {
	case "http":
		s, c, d, err := checkers.HTTPCheck(target.URL, target.ExpectedStatus, target.TimeoutSeconds)
		status = s
		code = &c
		duration = d
		if err != nil {
			e := err.Error()
			errMsg = &e
		}

	case "tcp":
		s, d, err := checkers.TCPCheck(target.URL, target.TimeoutSeconds)
		status = s
		duration = d
		if err != nil {
			e := err.Error()
			errMsg = &e
		}

	case "ssl":
		s, d, err := checkers.SSLCheck(target.URL, target.TimeoutSeconds)
		status = s
		duration = d
		if err != nil {
			e := err.Error()
			errMsg = &e
		}

	default:
		log.Println("Unknown check type:", target.Type)
		return
	}

	check := models.Check{
		TargetID:       target.ID,
		Status:         status,
		HttpCode:       code,
		ResponseTimeMs: duration,
		Error:          errMsg,
		Region:         target.RegionRestriction,
	}

	if err := w.checkRepo.Save(ctx, &check); err != nil {
		log.Println("Error saving check result:", err)
	}
}
