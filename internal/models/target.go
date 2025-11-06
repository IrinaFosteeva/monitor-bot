package models

import "time"

type Target struct {
	ID              int64     `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	URL             string    `db:"url" json:"url"`
	Method          string    `db:"method" json:"method"`
	ExpectedStatus  int       `db:"expected_status" json:"expected_status"`
	BodyRegex       *string   `db:"body_regex" json:"body_regex,omitempty"`
	IntervalSeconds int       `db:"interval_seconds" json:"interval_seconds"`
	TimeoutSeconds  int       `db:"timeout_seconds" json:"timeout_seconds"`
	Type            string    `db:"type" json:"type"` // "http", "tcp", "ssl"
	RegionID        int64     `db:"region_id" json:"region_id"`
	CreatedBy       *int64    `db:"created_by" json:"created_by,omitempty"`
	Enabled         bool      `db:"enabled" json:"enabled"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}
