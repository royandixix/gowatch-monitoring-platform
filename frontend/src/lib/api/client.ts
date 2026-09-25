import { env } from '$env/dynamic/public';
import { clearSession, getToken } from '$lib/utils/storage';

const API_BASE_URL = env.PUBLIC_API_BASE_URL ?? 'http://localhost:8081/api/v1';

interface ApiRequestOptions extends Omit<RequestInit, 'body'> {
	auth?: boolean;
	body?: unknown;
}

export class ApiError extends Error {
	status: number;

	constructor(message: string, status: number) {
		super(message);

		this.name = 'ApiError';
		this.status = status;
	}
}

function getErrorMessage(data: unknown): string {
	if (
		typeof data === 'object' &&
		data !== null &&
		'error' in data &&
		typeof data.error === 'string'
	) {
		return data.error;
	}

	if (
		typeof data === 'object' &&
		data !== null &&
		'message' in data &&
		typeof data.message === 'string'
	) {
		return data.message;
	}

	return 'Terjadi kesalahan pada server';
}

export async function apiRequest<T>(path: string, options: ApiRequestOptions = {}): Promise<T> {
	const { auth = false, body, headers: customHeaders, ...requestOptions } = options;

	const headers = new Headers(customHeaders);

	if (body !== undefined) {
		headers.set('Content-Type', 'application/json');
	}

	if (auth) {
		const token = getToken();

		if (!token) {
			clearSession();

			throw new ApiError('Sesi login tidak ditemukan', 401);
		}

		headers.set('Authorization', `Bearer ${token}`);
	}

	const response = await fetch(`${API_BASE_URL}${path}`, {
		...requestOptions,
		headers,
		body: body === undefined ? undefined : JSON.stringify(body)
	});

	if (response.status === 204) {
		return undefined as T;
	}

	const responseText = await response.text();

	let data: unknown = null;

	if (responseText) {
		try {
			data = JSON.parse(responseText);
		} catch {
			data = responseText;
		}
	}

	if (!response.ok) {
		if (response.status === 401 && auth) {
			clearSession();
		}

		throw new ApiError(getErrorMessage(data), response.status);
	}

	return data as T;
}
