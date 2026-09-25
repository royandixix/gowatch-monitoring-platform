export type HistoryStatus = 'ALL' | 'UP' | 'DOWN';

export interface MonitorResultRecord {
	id: number;
	monitor_id: number;
	monitor_name: string;
	monitor_url: string;
	status: string;
	http_status: number;
	response_time_ms: number;
	error: string;
	checked_at: string;
}

export interface HistoryQuery {
	status?: HistoryStatus;
	monitorId?: number;
	page?: number;
	perPage?: number;
}

export interface HistoryResponse {
	total: number;
	page: number;
	per_page: number;
	total_pages: number;
	results: MonitorResultRecord[];
}

export interface MonitorResultsResponse {
	total: number;
	results: MonitorResultRecord[];
}

export interface MonitorStats {
	monitor_id: number;
	total_checks: number;
	total_up: number;
	total_down: number;
	uptime_percentage: number;
	average_response_time_ms: number;
	last_response_time_ms: number;
	last_status: string;
	last_http_status: number;
	last_checked_at: string | null;
}

export interface MonitorStatsResponse {
	stats: MonitorStats;
}
