package model

import "time"

type MonitorResultRecord struct {
	ID             int64     `json:"id"`
	MonitorID      int64     `json:"monitor_id"`
	MonitorName    string    `json:"monitor_name"`
	MonitorURL     string    `json:"monitor_url"`
	Status         string    `json:"status"`
	HTTPStatus     int       `json:"http_status"`
	ResponseTimeMS int64     `json:"response_time_ms"`
	Error          string    `json:"error"`
	CheckedAt      time.Time `json:"checked_at"`
}

type MonitorHistoryFilter struct {
	Status    string
	MonitorID int64
	Page      int
	PerPage   int
}

type MonitorStats struct {
	MonitorID             int64      `json:"monitor_id"`
	TotalChecks           int64      `json:"total_checks"`
	TotalUp               int64      `json:"total_up"`
	TotalDown             int64      `json:"total_down"`
	UptimePercentage      float64    `json:"uptime_percentage"`
	AverageResponseTimeMS int64      `json:"average_response_time_ms"`
	LastResponseTimeMS    int64      `json:"last_response_time_ms"`
	LastStatus            string     `json:"last_status"`
	LastHTTPStatus        int        `json:"last_http_status"`
	LastCheckedAt         *time.Time `json:"last_checked_at"`
}
