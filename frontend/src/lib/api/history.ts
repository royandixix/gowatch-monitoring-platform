import { apiRequest } from '$lib/api/client';

import type {
	HistoryQuery,
	HistoryResponse,
	MonitorResultsResponse,
	MonitorStatsResponse
} from '$lib/types/history';

export function getHistory(options: HistoryQuery | number = {}): Promise<HistoryResponse> {
	const params = new URLSearchParams();

	if (typeof options === 'number') {
		params.set('per_page', String(options));
	} else {
		if (options.status && options.status !== 'ALL') {
			params.set('status', options.status);
		}

		if (options.monitorId && options.monitorId > 0) {
			params.set('monitor_id', String(options.monitorId));
		}

		if (options.page && options.page > 0) {
			params.set('page', String(options.page));
		}

		if (options.perPage && options.perPage > 0) {
			params.set('per_page', String(options.perPage));
		}
	}

	const query = params.toString();

	return apiRequest<HistoryResponse>(query ? `/history?${query}` : '/history', {
		method: 'GET',
		auth: true
	});
}

export function getMonitorResults(monitorId: number, limit = 100): Promise<MonitorResultsResponse> {
	return apiRequest<MonitorResultsResponse>(`/monitors/${monitorId}/results?limit=${limit}`, {
		method: 'GET',
		auth: true
	});
}

export function getMonitorStats(monitorId: number): Promise<MonitorStatsResponse> {
	return apiRequest<MonitorStatsResponse>(`/monitors/${monitorId}/stats`, {
		method: 'GET',
		auth: true
	});
}
