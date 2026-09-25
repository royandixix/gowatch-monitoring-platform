import { apiRequest } from '$lib/api/client';

import type {
	CreateMonitorPayload,
	DeleteMonitorResponse,
	MonitorListResponse,
	MonitorMutationResponse,
	MonitorResponse,
	UpdateMonitorPayload
} from '$lib/types/monitor';

export function getMonitors(): Promise<MonitorListResponse> {
	return apiRequest<MonitorListResponse>('/monitors', {
		method: 'GET',
		auth: true
	});
}

export function getMonitor(id: number): Promise<MonitorResponse> {
	return apiRequest<MonitorResponse>(`/monitors/${id}`, {
		method: 'GET',
		auth: true
	});
}

export function createMonitor(payload: CreateMonitorPayload): Promise<MonitorMutationResponse> {
	return apiRequest<MonitorMutationResponse>('/monitors', {
		method: 'POST',
		auth: true,
		body: payload
	});
}

export function updateMonitor(
	id: number,
	payload: UpdateMonitorPayload
): Promise<MonitorMutationResponse> {
	return apiRequest<MonitorMutationResponse>(`/monitors/${id}`, {
		method: 'PUT',
		auth: true,
		body: payload
	});
}

export function deleteMonitor(id: number): Promise<DeleteMonitorResponse> {
	return apiRequest<DeleteMonitorResponse>(`/monitors/${id}`, {
		method: 'DELETE',
		auth: true
	});
}
