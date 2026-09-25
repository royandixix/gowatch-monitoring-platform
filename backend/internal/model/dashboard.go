package model

import "time"

type DashboardSummary struct {
	TotalMonitors          int64   `json:"total_monitors"`
	UpMonitors             int64   `json:"up_monitors"`
	DownMonitors           int64   `json:"down_monitors"`
	PausedMonitors         int64   `json:"paused_monitors"`
	WaitingMonitors        int64   `json:"waiting_monitors"`
	TotalChecks24H         int64   `json:"total_checks_24h"`
	TotalChecks7D          int64   `json:"total_checks_7d"`
	Uptime24HPercentage    float64 `json:"uptime_24h_percentage"`
	Uptime7DPercentage     float64 `json:"uptime_7d_percentage"`
	AverageResponseTime24H int64   `json:"average_response_time_24h_ms"`
}

type DashboardMonitorStatus struct {
	ID                 int64      `json:"id"`
	Name               string     `json:"name"`
	URL                string     `json:"url"`
	MonitorType        string     `json:"monitor_type"`
	IsActive           bool       `json:"is_active"`
	CurrentStatus      string     `json:"current_status"`
	LastHTTPStatus     int        `json:"last_http_status"`
	LastResponseTimeMS int64      `json:"last_response_time_ms"`
	LastCheckedAt      *time.Time `json:"last_checked_at"`
}

type DashboardActivity struct {
	ID             int64     `json:"id"`
	MonitorID      int64     `json:"monitor_id"`
	MonitorName    string    `json:"monitor_name"`
	Status         string    `json:"status"`
	HTTPStatus     int       `json:"http_status"`
	ResponseTimeMS int64     `json:"response_time_ms"`
	Error          string    `json:"error"`
	CheckedAt      time.Time `json:"checked_at"`
}

type DashboardChartPoint struct {
	ID             int64     `json:"id"`
	MonitorID      int64     `json:"monitor_id"`
	MonitorName    string    `json:"monitor_name"`
	Status         string    `json:"status"`
	ResponseTimeMS int64     `json:"response_time_ms"`
	CheckedAt      time.Time `json:"checked_at"`
}

type DashboardResponse struct {
	Summary        DashboardSummary         `json:"summary"`
	Monitors       []DashboardMonitorStatus `json:"monitors"`
	RecentActivity []DashboardActivity      `json:"recent_activity"`
	ResponseChart  []DashboardChartPoint    `json:"response_chart"`
}
