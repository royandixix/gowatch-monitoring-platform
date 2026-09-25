package model

import "time"

type Monitor struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Name            string    `json:"name"`
	URL             string    `json:"url"`
	MonitorType     string    `json:"monitor_type"`
	IntervalSeconds int       `json:"interval_seconds"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateMonitorRequest struct {
	Name            string `json:"name" binding:"required,min=2,max=150"`
	URL             string `json:"url" binding:"required"`
	MonitorType     string `json:"monitor_type"`
	IntervalSeconds int    `json:"interval_seconds"`
	IsActive        *bool  `json:"is_active"`
}

type UpdateMonitorRequest struct {
	Name            *string `json:"name"`
	URL             *string `json:"url"`
	MonitorType     *string `json:"monitor_type"`
	IntervalSeconds *int    `json:"interval_seconds"`
	IsActive        *bool   `json:"is_active"`
}

type MonitorResult struct {
	URL          string `json:"url"`
	Status       string `json:"status"`
	HTTPStatus   int    `json:"http_status"`
	ResponseTime int64  `json:"response_time_ms"`
	CheckedAt    string `json:"checked_at"`
	Error        string `json:"error,omitempty"`
}

type MultipleCheckResponse struct {
	Total           int             `json:"total"`
	TotalDurationMs int64           `json:"total_duration_ms"`
	Results         []MonitorResult `json:"results"`
}
