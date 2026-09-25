import { apiRequest } from '$lib/api/client';

import type {
	ChangePasswordPayload,
	ChangePasswordResponse,
	LoginPayload,
	LoginResponse,
	MeResponse,
	RegisterPayload,
	RegisterResponse,
	UpdateProfilePayload,
	UpdateProfileResponse
} from '$lib/types/auth';

export function login(payload: LoginPayload): Promise<LoginResponse> {
	return apiRequest<LoginResponse>('/auth/login', {
		method: 'POST',
		body: payload
	});
}

export function register(payload: RegisterPayload): Promise<RegisterResponse> {
	return apiRequest<RegisterResponse>('/auth/register', {
		method: 'POST',
		body: payload
	});
}

export function getMe(): Promise<MeResponse> {
	return apiRequest<MeResponse>('/auth/me', {
		method: 'GET',
		auth: true
	});
}

export function updateProfile(payload: UpdateProfilePayload): Promise<UpdateProfileResponse> {
	return apiRequest<UpdateProfileResponse>('/profile', {
		method: 'PUT',
		auth: true,
		body: payload
	});
}

export function changePassword(payload: ChangePasswordPayload): Promise<ChangePasswordResponse> {
	return apiRequest<ChangePasswordResponse>('/profile/password', {
		method: 'PUT',
		auth: true,
		body: payload
	});
}
