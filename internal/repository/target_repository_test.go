package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"monitor-bot/internal/models"
)

func TestCreateTarget(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewTargetRepository(sqlxDB)
	ctx := context.Background()

	target := &models.Target{
		Title:       "Test",
		Description: "Desc",
		Deadline:    time.Now(),
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO targets (title, description, deadline, created_at, updated_at)
		 VALUES ($1, $2, $3, NOW(), NOW())
		 RETURNING id, created_at, updated_at`)).
		WithArgs(target.Title, target.Description, target.Deadline).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	err := repo.Create(ctx, target)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), target.ID)
}
