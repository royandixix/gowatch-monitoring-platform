import { writable } from 'svelte/store';

import { getMe } from '$lib/api/auth';

import type { AuthState, User } from '$lib/types/auth';

import { clearSession, getToken, saveSession } from '$lib/utils/storage';

const initialState: AuthState = {
	user: null,
	isAuthenticated: false,
	initialized: false
};

function createAuthStore() {
	const { subscribe, set } = writable<AuthState>(initialState);

	async function initialize(): Promise<void> {
		const token = getToken();

		if (!token) {
			set({
				user: null,
				isAuthenticated: false,
				initialized: true
			});

			return;
		}

		try {
			const response = await getMe();

			set({
				user: response.user,
				isAuthenticated: true,
				initialized: true
			});
		} catch {
			clearSession();

			set({
				user: null,
				isAuthenticated: false,
				initialized: true
			});
		}
	}

	function setAuthenticatedUser(token: string, user: User): void {
		saveSession(token, user);

		set({
			user,
			isAuthenticated: true,
			initialized: true
		});
	}

	function updateUser(user: User): void {
		const token = getToken();

		if (token) {
			saveSession(token, user);
		}

		set({
			user,
			isAuthenticated: true,
			initialized: true
		});
	}

	function logout(): void {
		clearSession();

		set({
			user: null,
			isAuthenticated: false,
			initialized: true
		});
	}

	return {
		subscribe,
		initialize,
		setAuthenticatedUser,
		updateUser,
		logout
	};
}

export const authStore = createAuthStore();
