<script lang="ts">
	import { goto } from '$app/navigation';
	import { base, resolve } from '$app/paths';

	import { fade } from 'svelte/transition';

	import { login, register } from '$lib/api/auth';

	import { ApiError } from '$lib/api/client';

	import { authStore } from '$lib/stores/auth';

	import { showError, showSuccess } from '$lib/utils/alert';

	type AuthMode = 'login' | 'register';

	interface Props {
		initialMode: AuthMode;
	}

	let { initialMode }: Props = $props();

	let mode = $derived(initialMode);

	let loading = $state(false);

	let switching = $state(false);

	let loginEmail = $state('');

	let loginPassword = $state('');

	let showLoginPassword = $state(false);

	let registerName = $state('');

	let registerEmail = $state('');

	let registerPassword = $state('');

	let registerConfirmPassword = $state('');

	let showRegisterPassword = $state(false);

	const illustrationUrl = `${base}/illustrations/login.svg`;

	async function switchMode(nextMode: AuthMode): Promise<void> {
		if (mode === nextMode || switching) {
			return;
		}

		switching = true;

		await new Promise<void>((resolveDelay) => {
			window.setTimeout(resolveDelay, 420);
		});

		try {
			await goto(resolve(nextMode === 'login' ? '/login' : '/register'));
		} finally {
			switching = false;
		}
	}

	async function handleLogin(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		const email = loginEmail.trim().toLowerCase();

		if (!email || !loginPassword) {
			await showError('Data belum lengkap', 'Email dan password wajib diisi.');

			return;
		}

		loading = true;

		try {
			const response = await login({
				email,
				password: loginPassword
			});

			authStore.setAuthenticatedUser(response.token, response.user);

			await goto(resolve('/dashboard'));
		} catch (error) {
			if (error instanceof ApiError) {
				await showError('Login gagal', error.message);
			} else {
				await showError('Server tidak terhubung', 'Pastikan backend GoWatch sedang berjalan.');
			}
		} finally {
			loading = false;
		}
	}

	async function handleRegister(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		const name = registerName.trim();

		const email = registerEmail.trim().toLowerCase();

		if (!name || !email || !registerPassword || !registerConfirmPassword) {
			await showError(
				'Data belum lengkap',
				'Nama, email, password, dan konfirmasi password wajib diisi.'
			);

			return;
		}

		if (name.length < 2) {
			await showError('Nama tidak valid', 'Nama minimal harus terdiri dari 2 karakter.');

			return;
		}

		if (registerPassword.length < 8) {
			await showError('Password terlalu pendek', 'Password minimal harus terdiri dari 8 karakter.');

			return;
		}

		if (registerPassword.length > 72) {
			await showError('Password terlalu panjang', 'Password maksimal terdiri dari 72 karakter.');

			return;
		}

		if (registerPassword !== registerConfirmPassword) {
			await showError('Password tidak sama', 'Password dan konfirmasi password harus sama.');

			return;
		}

		loading = true;

		try {
			await register({
				name,
				email,
				password: registerPassword
			});

			await showSuccess(
				'Registrasi berhasil',
				'Akun GoWatch berhasil dibuat. Silakan login menggunakan akun tersebut.'
			);

			registerName = '';
			registerEmail = '';
			registerPassword = '';
			registerConfirmPassword = '';

			await goto(resolve('/login'));
		} catch (error) {
			if (error instanceof ApiError) {
				await showError('Registrasi gagal', error.message);
			} else {
				await showError('Server tidak terhubung', 'Pastikan backend GoWatch sedang berjalan.');
			}
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>
		{mode === 'login' ? 'Login' : 'Register'} | GoWatch
	</title>

	<meta name="description" content="GoWatch Website & API Monitoring Platform" />
</svelte:head>

<div class="auth-page">
	<div class="desktop-background">
		<section class="visual visual-login">
			<div class="visual-content">
				<div class="brand">
					<div class="brand-logo">G</div>

					<div>
						<strong> GoWatch </strong>

						<span> Website & API Monitoring </span>
					</div>
				</div>

				<div class="visual-main">
					<div class="visual-illustration">
						<img src={illustrationUrl} alt="Ilustrasi monitoring GoWatch" />
					</div>

					<div class="visual-copy">
						<p class="visual-label">Monitoring Platform</p>

						<h1>Pantau layanan dalam satu tempat.</h1>

						<p>
							GoWatch memantau status website dan API secara otomatis, mencatat waktu respons, serta
							menyimpan riwayat setiap pengecekan agar kondisi layanan lebih mudah dipantau.
						</p>
					</div>
				</div>
			</div>
		</section>

		<section class="visual visual-register">
			<div class="visual-content">
				<div class="brand">
					<div class="brand-logo dark-logo">G</div>

					<div>
						<strong> GoWatch </strong>

						<span> Website & API Monitoring </span>
					</div>
				</div>

				<div class="visual-main">
					<div class="visual-illustration">
						<img src={illustrationUrl} alt="Ilustrasi penggunaan GoWatch" />
					</div>

					<div class="visual-copy">
						<p class="visual-label">Mulai Monitoring</p>

						<h1>Pantau website dan API milikmu.</h1>

						<p>
							Buat akun untuk menambahkan layanan, menentukan interval pengecekan, dan melihat hasil
							monitoring dari dashboard GoWatch.
						</p>
					</div>
				</div>
			</div>
		</section>
	</div>

	<section class:panel-register={mode === 'register'} class="auth-panel">
		<div class="mobile-brand">
			<div class="mobile-logo">G</div>

			<div>
				<strong> GoWatch </strong>

				<span> Website & API Monitoring </span>
			</div>
		</div>

		{#key mode}
			<div
				class="form-wrapper"
				in:fade={{
					duration: 200,
					delay: 100
				}}
			>
				{#if mode === 'login'}
					<div class="form-container">
						<div class="form-heading">
							<p class="section-label">Selamat datang kembali</p>

							<h2>Login</h2>

							<p class="form-description">
								Masuk ke akun GoWatch untuk melihat dan mengelola layanan yang sedang dipantau.
							</p>
						</div>

						<form onsubmit={handleLogin} class="auth-form">
							<div class="field">
								<label for="loginEmail"> Email </label>

								<input
									id="loginEmail"
									type="email"
									bind:value={loginEmail}
									autocomplete="email"
									placeholder="nama@email.com"
								/>
							</div>

							<div class="field">
								<label for="loginPassword"> Password </label>

								<div class="password-input">
									<input
										id="loginPassword"
										type={showLoginPassword ? 'text' : 'password'}
										bind:value={loginPassword}
										autocomplete="current-password"
										placeholder="Masukkan password"
									/>

									<button
										type="button"
										onclick={() => {
											showLoginPassword = !showLoginPassword;
										}}
									>
										{showLoginPassword ? 'Sembunyikan' : 'Lihat'}
									</button>
								</div>
							</div>

							<button type="submit" class="primary-button" disabled={loading}>
								{#if loading}
									<span class="spinner"></span>

									Memproses...
								{:else}
									Login
								{/if}
							</button>
						</form>

						<div class="switch-section">
							<span> Belum memiliki akun? </span>

							<button
								type="button"
								disabled={switching}
								onclick={() => {
									void switchMode('register');
								}}
							>
								Daftar
							</button>
						</div>
					</div>
				{:else}
					<div class="form-container">
						<div class="form-heading">
							<p class="section-label">Buat akun baru</p>

							<h2>Daftar</h2>

							<p class="form-description">
								Buat akun GoWatch untuk mulai menambahkan dan memantau website maupun API.
							</p>
						</div>

						<form onsubmit={handleRegister} class="auth-form">
							<div class="field">
								<label for="registerName"> Nama lengkap </label>

								<input
									id="registerName"
									type="text"
									bind:value={registerName}
									autocomplete="name"
									placeholder="Masukkan nama lengkap"
								/>
							</div>

							<div class="field">
								<label for="registerEmail"> Email </label>

								<input
									id="registerEmail"
									type="email"
									bind:value={registerEmail}
									autocomplete="email"
									placeholder="nama@email.com"
								/>
							</div>

							<div class="field">
								<label for="registerPassword"> Password </label>

								<div class="password-input">
									<input
										id="registerPassword"
										type={showRegisterPassword ? 'text' : 'password'}
										bind:value={registerPassword}
										autocomplete="new-password"
										placeholder="Minimal 8 karakter"
									/>

									<button
										type="button"
										onclick={() => {
											showRegisterPassword = !showRegisterPassword;
										}}
									>
										{showRegisterPassword ? 'Sembunyikan' : 'Lihat'}
									</button>
								</div>
							</div>

							<div class="field">
								<label for="registerConfirmPassword"> Konfirmasi password </label>

								<input
									id="registerConfirmPassword"
									type={showRegisterPassword ? 'text' : 'password'}
									bind:value={registerConfirmPassword}
									autocomplete="new-password"
									placeholder="Ulangi password"
								/>
							</div>

							<button type="submit" class="primary-button" disabled={loading}>
								{#if loading}
									<span class="spinner"></span>

									Membuat akun...
								{:else}
									Daftar
								{/if}
							</button>
						</form>

						<div class="switch-section">
							<span> Sudah memiliki akun? </span>

							<button
								type="button"
								disabled={switching}
								onclick={() => {
									void switchMode('login');
								}}
							>
								Login
							</button>
						</div>
					</div>
				{/if}
			</div>
		{/key}

		<footer class="auth-footer">© 2026 GoWatch</footer>
	</section>
</div>

<style>
	:global(body) {
		overflow-x: hidden;
	}

	.auth-page {
		position: relative;
		width: 100%;
		min-height: 100svh;
		overflow: hidden;
		background: #ffffff;
	}

	.desktop-background {
		position: absolute;
		inset: 0;
		display: grid;
		grid-template-columns: 1fr 1fr;
	}

	.visual {
		position: relative;
		min-height: 100svh;
		overflow: hidden;
	}

	.visual-login {
		background: linear-gradient(135deg, #13a8ed 0%, #078bdc 55%, #0879cf 100%);
	}

	.visual-register {
		background: linear-gradient(135deg, #1c2430 0%, #151c26 55%, #10151d 100%);
	}

	.visual-content {
		display: flex;
		flex-direction: column;
		width: 100%;
		min-height: 100svh;
		padding: clamp(28px, 4vw, 56px);
		color: #ffffff;
	}

	.brand {
		display: flex;
		align-items: center;
		gap: 12px;
		flex-shrink: 0;
		animation: content-enter 500ms ease-out both;
	}

	.brand-logo {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 44px;
		height: 44px;
		flex-shrink: 0;
		border-radius: 8px;
		background: #ffffff;
		color: #078bdc;
		font-size: 18px;
		font-weight: 800;
	}

	.dark-logo {
		color: #161d27;
	}

	.brand strong {
		display: block;
		font-size: 15px;
		font-weight: 700;
	}

	.brand span {
		display: block;
		margin-top: 3px;
		font-size: 10px;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		color: rgba(255, 255, 255, 0.65);
	}

	.visual-main {
		display: flex;
		flex: 1;
		flex-direction: column;
		align-items: flex-start;
		justify-content: center;
		width: 100%;
		max-width: 560px;
		padding: 40px 0;
	}

	.visual-illustration {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		margin-bottom: 26px;
		animation: content-enter 550ms ease-out 50ms both;
	}

	.visual-illustration img {
		display: block;
		width: min(100%, 330px);
		height: auto;
		max-height: 210px;
		object-fit: contain;
	}

	.visual-copy {
		max-width: 520px;
		animation: content-enter 550ms ease-out 100ms both;
	}

	.visual-label {
		margin: 0 0 12px;
		font-size: 10px;
		font-weight: 600;
		letter-spacing: 0.12em;
		text-transform: uppercase;
		color: rgba(255, 255, 255, 0.62);
	}

	.visual-copy h1 {
		max-width: 520px;
		margin: 0;
		font-size: clamp(32px, 3.5vw, 52px);
		font-weight: 650;
		line-height: 1.08;
		letter-spacing: -0.04em;
	}

	.visual-copy > p:last-child {
		max-width: 470px;
		margin: 18px 0 0;
		font-size: 14px;
		line-height: 1.7;
		color: rgba(255, 255, 255, 0.72);
	}

	.auth-panel {
		position: absolute;
		top: 0;
		bottom: 0;
		left: 50%;
		z-index: 10;
		display: flex;
		flex-direction: column;
		width: 50%;
		min-height: 100svh;
		background: #ffffff;
		transition:
			left 450ms cubic-bezier(0.4, 0, 0.2, 1),
			box-shadow 450ms ease;
		box-shadow: -18px 0 55px rgba(15, 23, 42, 0.1);
	}

	.auth-panel.panel-register {
		left: 0;
		box-shadow: 18px 0 55px rgba(0, 0, 0, 0.16);
	}

	.form-wrapper {
		display: flex;
		flex: 1;
		align-items: center;
		justify-content: center;
		width: 100%;
		padding: clamp(32px, 5vw, 64px);
	}

	.form-container {
		width: min(100%, 400px);
	}

	.form-heading {
		margin-bottom: 32px;
	}

	.section-label {
		margin: 0;
		color: #0b8ce8;
		font-size: 12px;
		font-weight: 600;
	}

	.form-heading h2 {
		margin: 8px 0 0;
		color: #111827;
		font-size: 38px;
		font-weight: 500;
		letter-spacing: -0.035em;
	}

	.form-description {
		max-width: 360px;
		margin: 14px 0 0;
		color: #7b8492;
		font-size: 14px;
		line-height: 1.7;
	}

	.auth-form {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.field label {
		display: block;
		margin-bottom: 7px;
		color: #4b5563;
		font-size: 13px;
		font-weight: 500;
	}

	.field input {
		width: 100%;
		height: 47px;
		padding: 0 14px;
		border: 1px solid #d5dae1;
		border-radius: 6px;
		background: #ffffff;
		color: #111827;
		font-size: 14px;
		outline: none;
		transition:
			border-color 160ms,
			box-shadow 160ms;
	}

	.field input::placeholder {
		color: #a0a7b2;
	}

	.field input:focus {
		border-color: #0b8ce8;
		box-shadow: 0 0 0 3px rgba(11, 140, 232, 0.1);
	}

	.password-input {
		position: relative;
	}

	.password-input input {
		padding-right: 100px;
	}

	.password-input button {
		position: absolute;
		top: 50%;
		right: 14px;
		transform: translateY(-50%);
		border: 0;
		background: transparent;
		color: #7c8491;
		font-size: 11px;
		font-weight: 600;
		cursor: pointer;
	}

	.password-input button:hover {
		color: #0b8ce8;
	}

	.primary-button {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 47px;
		margin-top: 2px;
		border: 0;
		border-radius: 6px;
		background: #0b8ce8;
		color: #ffffff;
		font-size: 13px;
		font-weight: 700;
		cursor: pointer;
		transition:
			background 160ms,
			box-shadow 160ms,
			transform 160ms;
	}

	.primary-button:hover {
		background: #087acb;
		box-shadow: 0 5px 14px rgba(11, 140, 232, 0.18);
	}

	.primary-button:active {
		transform: translateY(1px);
	}

	.primary-button:disabled {
		cursor: not-allowed;
		opacity: 0.65;
	}

	.switch-section {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 5px;
		margin-top: 28px;
		color: #7b8492;
		font-size: 13px;
	}

	.switch-section button {
		border: 0;
		background: transparent;
		color: #0b8ce8;
		font-size: 13px;
		font-weight: 700;
		cursor: pointer;
	}

	.switch-section button:hover {
		text-decoration: underline;
	}

	.switch-section button:disabled {
		cursor: default;
		opacity: 0.6;
	}

	.spinner {
		width: 14px;
		height: 14px;
		margin-right: 8px;
		border: 2px solid rgba(255, 255, 255, 0.4);
		border-top-color: #ffffff;
		border-radius: 50%;
		animation: spin 700ms linear infinite;
	}

	.auth-footer {
		padding: 0 36px 26px;
		color: #a0a7b2;
		font-size: 10px;
		text-align: right;
	}

	.mobile-brand {
		display: none;
	}

	.mobile-logo {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 40px;
		height: 40px;
		border-radius: 7px;
		background: #0b8ce8;
		color: #ffffff;
		font-size: 16px;
		font-weight: 800;
	}

	.mobile-brand strong {
		display: block;
		color: #111827;
		font-size: 14px;
	}

	.mobile-brand span {
		display: block;
		margin-top: 2px;
		color: #8b94a2;
		font-size: 10px;
		text-transform: uppercase;
		letter-spacing: 0.08em;
	}

	@keyframes content-enter {
		from {
			opacity: 0;
			transform: translateY(12px);
		}

		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}

	@media (max-height: 820px) and (min-width: 1024px) {
		.visual-content {
			padding-top: 24px;
			padding-bottom: 24px;
		}

		.visual-main {
			padding: 24px 0;
		}

		.visual-illustration {
			margin-bottom: 16px;
		}

		.visual-illustration img {
			max-height: 155px;
		}

		.visual-copy h1 {
			font-size: clamp(29px, 3vw, 43px);
		}

		.visual-copy > p:last-child {
			margin-top: 12px;
			font-size: 13px;
			line-height: 1.55;
		}
	}

	@media (max-width: 1023px) {
		.desktop-background {
			display: none;
		}

		.auth-page {
			overflow: visible;
			background: #ffffff;
		}

		.auth-panel,
		.auth-panel.panel-register {
			position: relative;
			left: 0;
			width: 100%;
			min-height: 100svh;
			box-shadow: none;
			transition: none;
		}

		.mobile-brand {
			display: flex;
			align-items: center;
			gap: 11px;
			padding: 24px 24px 0;
		}

		.form-wrapper {
			padding: 52px 24px 36px;
		}

		.auth-footer {
			padding: 0 24px 24px;
			text-align: center;
		}
	}

	@media (max-width: 640px) {
		.form-wrapper {
			align-items: flex-start;
			padding: 48px 20px 32px;
		}

		.form-heading {
			margin-bottom: 28px;
		}

		.form-heading h2 {
			font-size: 34px;
		}

		.form-description {
			font-size: 13px;
		}

		.auth-form {
			gap: 18px;
		}

		.field input,
		.primary-button {
			height: 48px;
		}

		.switch-section {
			margin-top: 24px;
		}
	}

	@media (max-width: 380px) {
		.mobile-brand {
			padding: 18px 16px 0;
		}

		.form-wrapper {
			padding: 40px 16px 28px;
		}

		.form-heading h2 {
			font-size: 30px;
		}

		.password-input input {
			padding-right: 86px;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.brand,
		.visual-illustration,
		.visual-copy,
		.spinner {
			animation: none;
		}

		.auth-panel {
			transition: none;
		}
	}
</style>
