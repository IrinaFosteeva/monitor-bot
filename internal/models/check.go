package models

import "time"

type Check struct {
	ID             int64     `db:"id" json:"id"`
	TargetID       int64     `db:"target_id" json:"target_id"`
	Timestamp      time.Time `db:"timestamp" json:"timestamp"`
	Status         string    `db:"status" json:"status"`
	HttpCode       *int      `db:"http_code" json:"http_code,omitempty"`
	ResponseTimeMs int64     `db:"response_time_ms" json:"response_time_ms"`
	Error          *string   `db:"error" json:"error,omitempty"`
	Region         *string   `db:"region" json:"region,omitempty"`
}
