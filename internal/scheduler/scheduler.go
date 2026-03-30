package scheduler

import (
	"context"
	"log"
	"time"

	"monitor-bot/internal/models"
	"monitor-bot/internal/repository"
)

type targetRunner interface {
	Run(ctx context.Context, target models.Target)
}

type Scheduler struct {
	targetRepo  repository.TargetRepositoryInterface
	runner      targetRunner
	interval    time.Duration
	now         func() time.Time
	lastStarted map[int64]time.Time
	running     map[int64]bool
}

func NewScheduler(tRepo repository.TargetRepositoryInterface, runner targetRunner, interval time.Duration) *Scheduler {
	return &Scheduler{
		targetRepo:  tRepo,
		runner:      runner,
		interval:    interval,
		now:         time.Now,
		lastStarted: make(map[int64]time.Time),
		running:     make(map[int64]bool),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	done := make(chan int64, 128)

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return
		case targetID := <-done:
			delete(s.running, targetID)
		case <-ticker.C:
			targets, err := s.targetRepo.GetAll(ctx)
			if err != nil {
				log.Println("Scheduler error:", err)
				continue
			}

			now := s.now()
			for _, t := range targets {
				if s.running[t.ID] {
					continue
				}
				if !s.shouldRun(now, t) {
					continue
				}

				s.running[t.ID] = true
				s.lastStarted[t.ID] = now

				go func(target models.Target) {
					defer func() {
						done <- target.ID
					}()
					s.runner.Run(ctx, target)
				}(t)
			}
		}
	}
}

func (s *Scheduler) shouldRun(now time.Time, target models.Target) bool {
	lastStarted, ok := s.lastStarted[target.ID]
	if !ok {
		return true
	}

	return now.Sub(lastStarted) >= targetInterval(target)
}

func targetInterval(target models.Target) time.Duration {
	if target.IntervalSeconds <= 0 {
		return time.Minute
	}
	return time.Duration(target.IntervalSeconds) * time.Second
}
