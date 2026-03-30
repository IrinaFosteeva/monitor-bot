package scheduler

import (
	"context"
	"sync"
	"testing"
	"time"

	"monitor-bot/internal/models"
)

type fakeTargetRepo struct {
	targets []models.Target
}

func (r *fakeTargetRepo) Create(ctx context.Context, t *models.Target) error { return nil }
func (r *fakeTargetRepo) GetByID(ctx context.Context, id int64) (*models.Target, error) {
	return nil, nil
}
func (r *fakeTargetRepo) GetByURL(ctx context.Context, url string) (*models.Target, error) {
	return nil, nil
}
func (r *fakeTargetRepo) GetAll(ctx context.Context) ([]models.Target, error) {
	return r.targets, nil
}
func (r *fakeTargetRepo) Update(ctx context.Context, t *models.Target) error { return nil }
func (r *fakeTargetRepo) Delete(ctx context.Context, id int64) error         { return nil }
func (r *fakeTargetRepo) Disable(ctx context.Context, id int64) error        { return nil }
func (r *fakeTargetRepo) Enable(ctx context.Context, id int64) error         { return nil }

type fakeRunner struct {
	mu      sync.Mutex
	calls   map[int64]int
	blockCh chan struct{}
}

func newFakeRunner() *fakeRunner {
	return &fakeRunner{
		calls: make(map[int64]int),
	}
}

func (r *fakeRunner) Run(ctx context.Context, target models.Target) {
	r.mu.Lock()
	r.calls[target.ID]++
	blockCh := r.blockCh
	r.mu.Unlock()

	if blockCh != nil {
		select {
		case <-blockCh:
		case <-ctx.Done():
		}
	}
}

func (r *fakeRunner) Calls(targetID int64) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.calls[targetID]
}

func TestSchedulerRespectsTargetInterval(t *testing.T) {
	repo := &fakeTargetRepo{
		targets: []models.Target{
			{ID: 1, IntervalSeconds: 1},
		},
	}
	runner := newFakeRunner()
	scheduler := NewScheduler(repo, runner, 10*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go scheduler.Start(ctx)

	time.Sleep(150 * time.Millisecond)
	if got := runner.Calls(1); got != 1 {
		t.Fatalf("expected first target to run once within interval window, got %d", got)
	}

	time.Sleep(1100 * time.Millisecond)
	cancel()

	if got := runner.Calls(1); got < 2 {
		t.Fatalf("expected target to run again after its interval elapsed, got %d", got)
	}
}

func TestSchedulerDoesNotStartSameTargetWhilePreviousRunIsActive(t *testing.T) {
	repo := &fakeTargetRepo{
		targets: []models.Target{
			{ID: 1, IntervalSeconds: 1},
		},
	}
	runner := newFakeRunner()
	runner.blockCh = make(chan struct{})
	scheduler := NewScheduler(repo, runner, 10*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go scheduler.Start(ctx)

	time.Sleep(120 * time.Millisecond)
	if got := runner.Calls(1); got != 1 {
		t.Fatalf("expected only one in-flight run for target, got %d", got)
	}

	close(runner.blockCh)
	time.Sleep(30 * time.Millisecond)
	cancel()
}
