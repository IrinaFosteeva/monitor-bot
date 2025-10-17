package scheduler

import (
	"context"
	"log"
	"time"

	"monitor-bot/internal/repository"
	"monitor-bot/internal/worker"
)

type Scheduler struct {
	targetRepo *repository.TargetRepository
	worker     *worker.Worker
	interval   time.Duration
}

func NewScheduler(tRepo *repository.TargetRepository, w *worker.Worker, interval time.Duration) *Scheduler {
	return &Scheduler{
		targetRepo: tRepo,
		worker:     w,
		interval:   interval,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return
		case <-ticker.C:
			targets, err := s.targetRepo.GetAll(ctx)
			if err != nil {
				log.Println("Scheduler error:", err)
				continue
			}

			for _, t := range targets {
				go s.worker.Run(ctx, t)
			}
		}
	}
}
