import { browser } from '$app/environment';
import type { User } from '$lib/types/auth';

const TOKEN_KEY = 'gowatch_token';
const USER_KEY = 'gowatch_user';

export function getToken(): string | null {
	if (!browser) {
		return null;
	}

	return localStorage.getItem(TOKEN_KEY);
}

export function getStoredUser(): User | null {
	if (!browser) {
		return null;
	}

	const value = localStorage.getItem(USER_KEY);

	if (!value) {
		return null;
	}

	try {
		return JSON.parse(value) as User;
	} catch {
		localStorage.removeItem(USER_KEY);

		return null;
	}
}

export function saveSession(token: string, user: User): void {
	if (!browser) {
		return;
	}

	localStorage.setItem(TOKEN_KEY, token);

	localStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearSession(): void {
	if (!browser) {
		return;
	}

	localStorage.removeItem(TOKEN_KEY);
	localStorage.removeItem(USER_KEY);
}
