export type MonitorType = 'website' | 'api';

export interface Monitor {
	id: number;
	user_id: number;
	name: string;
	url: string;
	monitor_type: MonitorType;
	interval_seconds: number;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateMonitorPayload {
	name: string;
	url: string;
	monitor_type?: MonitorType;
	interval_seconds?: number;
	is_active?: boolean;
}

export interface UpdateMonitorPayload {
	name?: string;
	url?: string;
	monitor_type?: MonitorType;
	interval_seconds?: number;
	is_active?: boolean;
}

export interface MonitorListResponse {
	total: number;
	monitors: Monitor[];
}

export interface MonitorResponse {
	monitor: Monitor;
}

export interface MonitorMutationResponse {
	message: string;
	monitor: Monitor;
}

export interface DeleteMonitorResponse {
	message: string;
}
