import { apiRequest } from '$lib/api/client';

import type { DashboardResponse } from '$lib/types/dashboard';

export function getDashboard(): Promise<DashboardResponse> {
	return apiRequest<DashboardResponse>('/dashboard', {
		method: 'GET',
		auth: true
	});
}
