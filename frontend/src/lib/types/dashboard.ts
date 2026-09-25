export interface DashboardSummary {
	total_monitors: number;
	up_monitors: number;
	down_monitors: number;
	paused_monitors: number;
	waiting_monitors: number;

	total_checks_24h: number;
	total_checks_7d: number;

	uptime_24h_percentage: number;
	uptime_7d_percentage: number;

	average_response_time_24h_ms: number;
}

export interface DashboardMonitorStatus {
	id: number;
	name: string;
	url: string;
	monitor_type: string;
	is_active: boolean;

	current_status: string;

	last_http_status: number;
	last_response_time_ms: number;

	last_checked_at: string | null;
}

export interface DashboardActivity {
	id: number;
	monitor_id: number;
	monitor_name: string;

	status: string;

	http_status: number;
	response_time_ms: number;

	error: string;

	checked_at: string;
}

export interface DashboardChartPoint {
	id: number;
	monitor_id: number;
	monitor_name: string;

	status: string;

	response_time_ms: number;
	checked_at: string;
}

export interface DashboardResponse {
	summary: DashboardSummary;

	monitors: DashboardMonitorStatus[];

	recent_activity: DashboardActivity[];

	response_chart: DashboardChartPoint[];
}
