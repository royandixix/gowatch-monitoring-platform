export interface User {
	id: number;
	name: string;
	email: string;
	created_at: string;
	updated_at: string;
}

export interface LoginPayload {
	email: string;
	password: string;
}

export interface RegisterPayload {
	name: string;
	email: string;
	password: string;
}

export interface UpdateProfilePayload {
	name: string;
	email: string;
}

export interface ChangePasswordPayload {
	current_password: string;
	new_password: string;
}

export interface LoginResponse {
	message: string;
	token: string;
	user: User;
}

export interface RegisterResponse {
	message: string;
	user: User;
}

export interface MeResponse {
	user: User;
}

export interface UpdateProfileResponse {
	message: string;
	user: User;
}

export interface ChangePasswordResponse {
	message: string;
}

export interface AuthState {
	user: User | null;
	isAuthenticated: boolean;
	initialized: boolean;
}
