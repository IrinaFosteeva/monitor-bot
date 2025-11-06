package repository

import (
	"context"
	"monitor-bot/internal/models"
)

type TargetRepositoryInterface interface {
	Create(ctx context.Context, t *models.Target) error
	GetByID(ctx context.Context, id int64) (*models.Target, error)
	GetByURL(ctx context.Context, url string) (*models.Target, error)
	GetAll(ctx context.Context) ([]models.Target, error)
	Update(ctx context.Context, t *models.Target) error
	Delete(ctx context.Context, id int64) error
}
