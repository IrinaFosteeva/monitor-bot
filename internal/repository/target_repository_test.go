package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"monitor-bot/internal/models"
)

func setupMockRepo(t *testing.T) (*TargetRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewTargetRepository(sqlxDB)
	return repo, mock
}

func TestCreateTarget_Success(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	target := &models.Target{
		Name:              "Example Site",
		URL:               "https://example.com",
		Method:            "GET",
		ExpectedStatus:    200,
		BodyRegex:         nil,
		IntervalSeconds:   60,
		TimeoutSeconds:    5,
		RegionRestriction: nil,
		CreatedBy:         nil,
		Enabled:           true,
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO targets
			(name, url, method, expected_status, body_regex, interval_seconds, timeout_seconds, region_restriction, created_by, enabled)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		 RETURNING id, created_at`)).
		WithArgs(
			target.Name,
			target.URL,
			target.Method,
			target.ExpectedStatus,
			target.BodyRegex,
			target.IntervalSeconds,
			target.TimeoutSeconds,
			target.RegionRestriction,
			target.CreatedBy,
			target.Enabled,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow(1, time.Now()))

	err := repo.Create(ctx, target)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), target.ID)
}

func TestCreateTarget_DBError(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	target := &models.Target{
		Name:              "Example Site",
		URL:               "https://example.com",
		Method:            "GET",
		ExpectedStatus:    200,
		BodyRegex:         nil,
		IntervalSeconds:   60,
		TimeoutSeconds:    5,
		RegionRestriction: nil,
		CreatedBy:         nil,
		Enabled:           true,
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO targets
			(name, url, method, expected_status, body_regex, interval_seconds, timeout_seconds, region_restriction, created_by, enabled)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		 RETURNING id, created_at`)).
		WithArgs(
			target.Name,
			target.URL,
			target.Method,
			target.ExpectedStatus,
			target.BodyRegex,
			target.IntervalSeconds,
			target.TimeoutSeconds,
			target.RegionRestriction,
			target.CreatedBy,
			target.Enabled,
		).
		WillReturnError(errors.New("db error"))

	err := repo.Create(ctx, target)
	assert.Error(t, err)
}

func TestGetByID_Success(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	expected := &models.Target{
		ID:                1,
		Name:              "Example Site",
		URL:               "https://example.com",
		Method:            "GET",
		ExpectedStatus:    200,
		BodyRegex:         nil,
		IntervalSeconds:   60,
		TimeoutSeconds:    5,
		RegionRestriction: nil,
		CreatedBy:         nil,
		Enabled:           true,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM targets WHERE id = $1")).
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "url", "method", "expected_status", "body_regex",
			"interval_seconds", "timeout_seconds", "region_restriction", "created_by", "enabled", "created_at",
		}).AddRow(
			expected.ID,
			expected.Name,
			expected.URL,
			expected.Method,
			expected.ExpectedStatus,
			expected.BodyRegex,
			expected.IntervalSeconds,
			expected.TimeoutSeconds,
			expected.RegionRestriction,
			expected.CreatedBy,
			expected.Enabled,
			time.Now(),
		))

	result, err := repo.GetByID(ctx, expected.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

func TestGetByID_NotFound(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM targets WHERE id = $1")).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetByID(ctx, 999)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, result)
}

func TestGetAll_Success(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM targets ORDER BY id")).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "url", "method", "expected_status", "body_regex",
			"interval_seconds", "timeout_seconds", "region_restriction", "created_by", "enabled", "created_at",
		}).AddRow(
			1, "Site1", "https://site1.com", "GET", 200, nil, 60, 5, nil, nil, true, time.Now(),
		).AddRow(
			2, "Site2", "https://site2.com", "HEAD", 200, nil, 120, 5, nil, nil, true, time.Now(),
		))

	result, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestGetAll_DBError(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM targets ORDER BY id")).
		WillReturnError(errors.New("db error"))

	result, err := repo.GetAll(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUpdate_Success(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	target := &models.Target{
		ID:                1,
		Name:              "New Site",
		URL:               "https://newsite.com",
		Method:            "POST",
		ExpectedStatus:    201,
		BodyRegex:         nil,
		IntervalSeconds:   30,
		TimeoutSeconds:    10,
		RegionRestriction: nil,
		Enabled:           true,
	}

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE targets SET
			name=$1,
			url=$2,
			method=$3,
			expected_status=$4,
			body_regex=$5,
			interval_seconds=$6,
			timeout_seconds=$7,
			region_restriction=$8,
			enabled=$9
		WHERE id=$10`)).
		WithArgs(
			target.Name,
			target.URL,
			target.Method,
			target.ExpectedStatus,
			target.BodyRegex,
			target.IntervalSeconds,
			target.TimeoutSeconds,
			target.RegionRestriction,
			target.Enabled,
			target.ID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(ctx, target)
	assert.NoError(t, err)
}

func TestUpdate_DBError(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	target := &models.Target{
		ID:                1,
		Name:              "New Site",
		URL:               "https://newsite.com",
		Method:            "POST",
		ExpectedStatus:    201,
		BodyRegex:         nil,
		IntervalSeconds:   30,
		TimeoutSeconds:    10,
		RegionRestriction: nil,
		Enabled:           true,
	}

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE targets SET
			name=$1,
			url=$2,
			method=$3,
			expected_status=$4,
			body_regex=$5,
			interval_seconds=$6,
			timeout_seconds=$7,
			region_restriction=$8,
			enabled=$9
		WHERE id=$10`)).
		WithArgs(
			target.Name,
			target.URL,
			target.Method,
			target.ExpectedStatus,
			target.BodyRegex,
			target.IntervalSeconds,
			target.TimeoutSeconds,
			target.RegionRestriction,
			target.Enabled,
			target.ID,
		).
		WillReturnError(errors.New("db error"))

	err := repo.Update(ctx, target)
	assert.Error(t, err)
}

func TestDelete_Success(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM targets WHERE id=$1")).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestDelete_DBError(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM targets WHERE id=$1")).
		WithArgs(int64(1)).
		WillReturnError(errors.New("db error"))

	err := repo.Delete(ctx, 1)
	assert.Error(t, err)
}
