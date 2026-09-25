<script lang="ts">
	import { onMount } from 'svelte';

	import { changePassword, updateProfile } from '$lib/api/auth';

	import PageHeader from '$lib/components/layout/PageHeader.svelte';

	import Button from '$lib/components/ui/Button.svelte';

	import Input from '$lib/components/ui/Input.svelte';

	import { authStore } from '$lib/stores/auth';

	import { showError, showSuccess } from '$lib/utils/alert';

	import { formatDate } from '$lib/utils/format';

	let name = $state('');

	let email = $state('');

	let profileLoading = $state(false);

	let currentPassword = $state('');

	let newPassword = $state('');

	let confirmPassword = $state('');

	let passwordLoading = $state(false);

	onMount(() => {
		syncProfileForm();
	});

	function syncProfileForm(): void {
		name = $authStore.user?.name ?? '';

		email = $authStore.user?.email ?? '';
	}

	function getInitial(value?: string): string {
		if (!value) {
			return 'U';
		}

		return value.trim().charAt(0).toUpperCase();
	}

	function getErrorMessage(error: unknown, fallback: string): string {
		if (error instanceof Error) {
			return error.message;
		}

		return fallback;
	}

	async function handleProfileSubmit(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		const cleanName = name.trim();

		const cleanEmail = email.trim().toLowerCase();

		if (cleanName.length < 2) {
			await showError('Nama tidak valid', 'Nama minimal harus terdiri dari 2 karakter.');

			return;
		}

		if (!cleanEmail) {
			await showError('Email belum diisi', 'Email wajib diisi.');

			return;
		}

		profileLoading = true;

		try {
			const response = await updateProfile({
				name: cleanName,
				email: cleanEmail
			});

			authStore.updateUser(response.user);

			name = response.user.name;

			email = response.user.email;

			await showSuccess('Profile diperbarui', 'Nama dan email akun berhasil diperbarui.');
		} catch (error) {
			await showError(
				'Gagal memperbarui profile',
				getErrorMessage(error, 'Terjadi kesalahan saat memperbarui profile.')
			);
		} finally {
			profileLoading = false;
		}
	}

	async function handlePasswordSubmit(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		if (!currentPassword || !newPassword || !confirmPassword) {
			await showError(
				'Data belum lengkap',
				'Password lama, password baru, dan konfirmasi password wajib diisi.'
			);

			return;
		}

		if (newPassword.length < 8) {
			await showError(
				'Password terlalu pendek',
				'Password baru minimal harus terdiri dari 8 karakter.'
			);

			return;
		}

		if (newPassword.length > 72) {
			await showError('Password terlalu panjang', 'Password maksimal terdiri dari 72 karakter.');

			return;
		}

		if (newPassword !== confirmPassword) {
			await showError('Password tidak sama', 'Password baru dan konfirmasi password harus sama.');

			return;
		}

		if (currentPassword === newPassword) {
			await showError('Password tidak berubah', 'Password baru harus berbeda dari password lama.');

			return;
		}

		passwordLoading = true;

		try {
			await changePassword({
				current_password: currentPassword,
				new_password: newPassword
			});

			currentPassword = '';
			newPassword = '';
			confirmPassword = '';

			await showSuccess('Password diperbarui', 'Password akun berhasil diubah.');
		} catch (error) {
			await showError(
				'Gagal mengubah password',
				getErrorMessage(error, 'Terjadi kesalahan saat mengubah password.')
			);
		} finally {
			passwordLoading = false;
		}
	}
</script>

<svelte:head>
	<title>Profile | GoWatch</title>
</svelte:head>

<PageHeader title="Profile" description="Kelola informasi akun dan keamanan GoWatch." />

<div class="-mt-14 px-4 pb-10 sm:px-6 lg:px-8">
	<div class="mx-auto max-w-6xl">
		<div class="overflow-hidden rounded-lg bg-white shadow-xl">
			<div class="h-32 bg-gradient-to-r from-slate-700 via-slate-800 to-slate-900"></div>

			<div class="px-6 pb-7 sm:px-8">
				<div class="-mt-12 flex flex-col gap-5 sm:flex-row sm:items-end sm:justify-between">
					<div class="flex items-end gap-5">
						<div
							class="flex h-24 w-24 shrink-0 items-center justify-center rounded-full border-4 border-white bg-sky-500 text-3xl font-bold text-white shadow-lg"
						>
							{getInitial($authStore.user?.name)}
						</div>

						<div class="min-w-0 pb-2">
							<h1 class="truncate text-xl font-semibold text-slate-800">
								{$authStore.user?.name ?? 'User'}
							</h1>

							<p class="mt-1 truncate text-sm text-slate-400">
								{$authStore.user?.email ?? ''}
							</p>
						</div>
					</div>

					<div class="pb-2">
						<span
							class="inline-flex items-center gap-2 rounded-full bg-emerald-50 px-3 py-1.5 text-xs font-semibold text-emerald-600"
						>
							<span class="h-2 w-2 rounded-full bg-emerald-500"></span>

							Active Account
						</span>
					</div>
				</div>

				<div class="mt-8 grid gap-4 border-t border-slate-100 pt-7 sm:grid-cols-3">
					<div class="rounded-lg border border-slate-100 bg-slate-50 p-4">
						<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">User ID</p>

						<p class="mt-2 text-sm font-semibold text-slate-700">
							#{$authStore.user?.id ?? '-'}
						</p>
					</div>

					<div class="rounded-lg border border-slate-100 bg-slate-50 p-4">
						<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">Member Sejak</p>

						<p class="mt-2 text-sm font-semibold text-slate-700">
							{$authStore.user?.created_at ? formatDate($authStore.user.created_at) : '-'}
						</p>
					</div>

					<div class="rounded-lg border border-slate-100 bg-slate-50 p-4">
						<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">Status</p>

						<p class="mt-2 flex items-center gap-2 text-sm font-semibold text-emerald-600">
							<span class="h-2 w-2 rounded-full bg-emerald-500"></span>

							Active
						</p>
					</div>
				</div>
			</div>
		</div>

		<div class="mt-6 grid gap-6 lg:grid-cols-2">
			<div class="rounded-lg bg-white p-6 shadow-lg sm:p-7">
				<div class="mb-6 flex items-center gap-4 border-b border-slate-100 pb-5">
					<div class="flex h-11 w-11 items-center justify-center rounded-lg bg-sky-50 text-sky-500">
						<svg
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="1.8"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="h-5 w-5"
							aria-hidden="true"
						>
							<path d="M20 21a8 8 0 0 0-16 0"></path>

							<circle cx="12" cy="7" r="4"></circle>
						</svg>
					</div>

					<div>
						<h2 class="text-sm font-semibold text-slate-800">Informasi Profile</h2>

						<p class="mt-1 text-xs text-slate-400">Perbarui nama dan email akun.</p>
					</div>
				</div>

				<form onsubmit={handleProfileSubmit} class="space-y-5">
					<Input
						id="profile-name"
						label="Nama lengkap"
						placeholder="Masukkan nama lengkap"
						autocomplete="name"
						required
						bind:value={name}
					/>

					<Input
						id="profile-email"
						label="Email"
						type="email"
						placeholder="nama@email.com"
						autocomplete="email"
						required
						bind:value={email}
					/>

					<div class="flex justify-end border-t border-slate-100 pt-5">
						<Button type="submit" loading={profileLoading}>
							<svg
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="h-4 w-4"
								aria-hidden="true"
							>
								<path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2Z"></path>

								<path d="M17 21v-8H7v8"></path>

								<path d="M7 3v5h8"></path>
							</svg>

							Simpan Profile
						</Button>
					</div>
				</form>
			</div>

			<div class="rounded-lg bg-white p-6 shadow-lg sm:p-7">
				<div class="mb-6 flex items-center gap-4 border-b border-slate-100 pb-5">
					<div
						class="flex h-11 w-11 items-center justify-center rounded-lg bg-amber-50 text-amber-500"
					>
						<svg
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="1.8"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="h-5 w-5"
							aria-hidden="true"
						>
							<rect x="5" y="11" width="14" height="10" rx="2"></rect>

							<path d="M8 11V7a4 4 0 0 1 8 0v4"></path>
						</svg>
					</div>

					<div>
						<h2 class="text-sm font-semibold text-slate-800">Keamanan Password</h2>

						<p class="mt-1 text-xs text-slate-400">Ganti password akun secara aman.</p>
					</div>
				</div>

				<form onsubmit={handlePasswordSubmit} class="space-y-5">
					<Input
						id="current-password"
						label="Password lama"
						type="password"
						placeholder="Masukkan password lama"
						autocomplete="current-password"
						required
						bind:value={currentPassword}
					/>

					<Input
						id="new-password"
						label="Password baru"
						type="password"
						placeholder="Minimal 8 karakter"
						autocomplete="new-password"
						required
						hint="Gunakan minimal 8 karakter dan jangan gunakan password lama."
						bind:value={newPassword}
					/>

					<Input
						id="confirm-password"
						label="Konfirmasi password baru"
						type="password"
						placeholder="Ulangi password baru"
						autocomplete="new-password"
						required
						bind:value={confirmPassword}
					/>

					<div class="flex justify-end border-t border-slate-100 pt-5">
						<Button type="submit" loading={passwordLoading}>
							<svg
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="h-4 w-4"
								aria-hidden="true"
							>
								<path d="M12 17v.01"></path>

								<rect x="5" y="10" width="14" height="11" rx="2"></rect>

								<path d="M8 10V7a4 4 0 0 1 8 0v3"></path>
							</svg>

							Ubah Password
						</Button>
					</div>
				</form>
			</div>
		</div>

		<div class="mt-6 rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
			<div class="flex items-start gap-4">
				<div
					class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-slate-100 text-slate-500"
				>
					<svg
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.8"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="h-4 w-4"
						aria-hidden="true"
					>
						<circle cx="12" cy="12" r="9"></circle>

						<path d="M12 11v5"></path>

						<path d="M12 8h.01"></path>
					</svg>
				</div>

				<div>
					<p class="text-sm font-semibold text-slate-700">Keamanan akun</p>

					<p class="mt-1 text-xs leading-6 text-slate-500">
						Password disimpan dalam bentuk hash bcrypt. GoWatch tidak menyimpan password asli di
						database. Saat mengganti password, password lama akan diverifikasi sebelum password baru
						disimpan.
					</p>
				</div>
			</div>
		</div>
	</div>
</div>
