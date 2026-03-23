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

func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTarget_Success(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	target := &models.Target{
		Name:            "Example Site",
		URL:             "https://example.com",
		Method:          "GET",
		ExpectedStatus:  200,
		BodyRegex:       nil,
		IntervalSeconds: 60,
		TimeoutSeconds:  5,
		RegionID:        1,
		Enabled:         true,
		Type:            "http",
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO targets
			(name, url, method, expected_status, body_regex, interval_seconds, timeout_seconds, region_id, enabled, type)
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
			target.RegionID,
			target.Enabled,
			target.Type,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow(1, time.Now()))

	err := repo.Create(ctx, target)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), target.ID)
	assertMockExpectations(t, mock)
}

func TestCreateTarget_AppliesDefaults(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	target := &models.Target{
		Name:            "Defaults Site",
		URL:             "https://defaults.example.com",
		Method:          "GET",
		ExpectedStatus:  200,
		IntervalSeconds: 60,
		TimeoutSeconds:  5,
		Enabled:         true,
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO targets
			(name, url, method, expected_status, body_regex, interval_seconds, timeout_seconds, region_id, enabled, type)
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
			int64(1),
			target.Enabled,
			"http",
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow(10, time.Now()))

	err := repo.Create(ctx, target)
	assert.NoError(t, err)
	assert.Equal(t, "http", target.Type)
	assert.Equal(t, int64(1), target.RegionID)
	assertMockExpectations(t, mock)
}

func TestGetByID_Success(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	expected := &models.Target{
		ID:              1,
		RegionID:        1,
		Name:            "Example Site",
		URL:             "https://example.com",
		Method:          "GET",
		ExpectedStatus:  200,
		BodyRegex:       nil,
		IntervalSeconds: 60,
		TimeoutSeconds:  5,
		Enabled:         true,
		Type:            "http",
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM targets WHERE id = $1")).
		WithArgs(expected.ID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "url", "method", "expected_status", "body_regex",
			"interval_seconds", "timeout_seconds", "region_id", "enabled", "created_at", "type",
		}).AddRow(
			expected.ID,
			expected.Name,
			expected.URL,
			expected.Method,
			expected.ExpectedStatus,
			expected.BodyRegex,
			expected.IntervalSeconds,
			expected.TimeoutSeconds,
			expected.RegionID,
			expected.Enabled,
			time.Now(),
			expected.Type,
		))

	result, err := repo.GetByID(ctx, expected.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assertMockExpectations(t, mock)
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
	assertMockExpectations(t, mock)
}

func TestGetAll_Success(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM targets WHERE enabled = TRUE ORDER BY id")).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "url", "method", "expected_status", "body_regex",
			"interval_seconds", "timeout_seconds", "region_id", "enabled", "created_at", "type",
		}).AddRow(
			1, "Site1", "https://site1.com", "GET", 200, nil, 60, 5, 1, true, time.Now(), "tcp",
		).AddRow(
			2, "Site2", "https://site2.com", "HEAD", 200, nil, 120, 5, 1, true, time.Now(), "tcp",
		))

	result, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assertMockExpectations(t, mock)
}

func TestGetAll_DBError(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM targets WHERE enabled = TRUE ORDER BY id")).
		WillReturnError(errors.New("db error"))

	result, err := repo.GetAll(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)
	assertMockExpectations(t, mock)
}

func TestUpdate_Success(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	target := &models.Target{
		ID:              1,
		Name:            "New Site",
		URL:             "https://newsite.com",
		Method:          "POST",
		ExpectedStatus:  201,
		BodyRegex:       nil,
		IntervalSeconds: 30,
		TimeoutSeconds:  10,
		RegionID:        1,
		Enabled:         true,
		Type:            "ssl",
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE targets SET
			name=$1,
			url=$2,
			method=$3,
			expected_status=$4,
			body_regex=$5,
			interval_seconds=$6,
			timeout_seconds=$7,
			region_id=$8,
			enabled=$9,
		    type=$10
		WHERE id=$11
		RETURNING created_at`)).
		WithArgs(
			target.Name,
			target.URL,
			target.Method,
			target.ExpectedStatus,
			target.BodyRegex,
			target.IntervalSeconds,
			target.TimeoutSeconds,
			target.RegionID,
			target.Enabled,
			target.Type,
			target.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(time.Now()))

	err := repo.Update(ctx, target)
	assert.NoError(t, err)
	assertMockExpectations(t, mock)
}

func TestUpdate_AppliesDefaults(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	target := &models.Target{
		ID:              1,
		Name:            "Updated Site",
		URL:             "https://updated.example.com",
		Method:          "GET",
		ExpectedStatus:  200,
		IntervalSeconds: 30,
		TimeoutSeconds:  10,
		Enabled:         false,
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE targets SET
			name=$1,
			url=$2,
			method=$3,
			expected_status=$4,
			body_regex=$5,
			interval_seconds=$6,
			timeout_seconds=$7,
			region_id=$8,
			enabled=$9,
		    type=$10
		WHERE id=$11
		RETURNING created_at`)).
		WithArgs(
			target.Name,
			target.URL,
			target.Method,
			target.ExpectedStatus,
			target.BodyRegex,
			target.IntervalSeconds,
			target.TimeoutSeconds,
			int64(1),
			target.Enabled,
			"http",
			target.ID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(time.Now()))

	err := repo.Update(ctx, target)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), target.RegionID)
	assert.Equal(t, "http", target.Type)
	assertMockExpectations(t, mock)
}

func TestUpdate_DBError(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	target := &models.Target{
		ID:              1,
		Name:            "New Site",
		URL:             "https://newsite.com",
		Method:          "POST",
		ExpectedStatus:  201,
		BodyRegex:       nil,
		IntervalSeconds: 30,
		TimeoutSeconds:  10,
		RegionID:        1,
		Enabled:         true,
		Type:            "ssl",
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE targets SET
			name=$1,
			url=$2,
			method=$3,
			expected_status=$4,
			body_regex=$5,
			interval_seconds=$6,
			timeout_seconds=$7,
			region_id=$8,
			enabled=$9,
		    type=$10
		WHERE id=$11
		RETURNING created_at`)).
		WithArgs(
			target.Name,
			target.URL,
			target.Method,
			target.ExpectedStatus,
			target.BodyRegex,
			target.IntervalSeconds,
			target.TimeoutSeconds,
			target.RegionID,
			target.Enabled,
			target.Type,
			target.ID,
		).
		WillReturnError(errors.New("db error"))

	err := repo.Update(ctx, target)
	assert.Error(t, err)
	assertMockExpectations(t, mock)
}

func TestUpdate_NotFound(t *testing.T) {
	repo, mock := setupMockRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	target := &models.Target{
		ID:              999,
		Name:            "Missing Site",
		URL:             "https://missing.example.com",
		Method:          "GET",
		ExpectedStatus:  200,
		IntervalSeconds: 30,
		TimeoutSeconds:  10,
		RegionID:        1,
		Enabled:         true,
		Type:            "http",
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE targets SET
			name=$1,
			url=$2,
			method=$3,
			expected_status=$4,
			body_regex=$5,
			interval_seconds=$6,
			timeout_seconds=$7,
			region_id=$8,
			enabled=$9,
		    type=$10
		WHERE id=$11
		RETURNING created_at`)).
		WithArgs(
			target.Name,
			target.URL,
			target.Method,
			target.ExpectedStatus,
			target.BodyRegex,
			target.IntervalSeconds,
			target.TimeoutSeconds,
			target.RegionID,
			target.Enabled,
			target.Type,
			target.ID,
		).
		WillReturnError(sql.ErrNoRows)

	err := repo.Update(ctx, target)
	assert.ErrorIs(t, err, ErrNotFound)
	assertMockExpectations(t, mock)
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
	assertMockExpectations(t, mock)
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
	assertMockExpectations(t, mock)
}
