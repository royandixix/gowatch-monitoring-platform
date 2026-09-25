<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { authStore } from '$lib/stores/auth';

	let {
		onMenuClick
	}: {
		onMenuClick: () => void;
	} = $props();

	let search = $state('');
	let accountOpen = $state(false);

	const title = $derived.by(() => {
		const pathname = page.url.pathname;

		if (pathname === '/dashboard') {
			return 'DASHBOARD';
		}

		if (pathname === '/monitors') {
			return 'MONITORS';
		}

		if (pathname.startsWith('/monitors/')) {
			return 'MONITOR DETAIL';
		}

		if (pathname === '/history') {
			return 'HISTORY';
		}

		if (pathname === '/profile') {
			return 'PROFILE';
		}

		return 'GOWATCH';
	});

	function getInitial(name?: string): string {
		if (!name) {
			return 'U';
		}

		return name.trim().charAt(0).toUpperCase();
	}

	async function handleSearch(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		const value = search.trim();

		if (!value) {
			await goto(resolve('/monitors'));
			return;
		}

		await goto(resolve(`/monitors?search=${encodeURIComponent(value)}`));
	}

	async function logout(): Promise<void> {
		accountOpen = false;
		authStore.logout();
		await goto(resolve('/login'));
	}
</script>

<header class="sticky top-0 z-30 h-16 bg-slate-800 text-white">
	<div class="flex h-full items-center justify-between px-4 sm:px-6 lg:px-8">
		<div class="flex items-center gap-3">
			<button
				type="button"
				aria-label="Buka sidebar"
				class="flex h-10 w-10 items-center justify-center rounded-md text-white/80 transition hover:bg-white/10 lg:hidden"
				onclick={onMenuClick}
			>
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="h-5 w-5">
					<path d="M4 6h16"></path>
					<path d="M4 12h16"></path>
					<path d="M4 18h16"></path>
				</svg>
			</button>

			<p class="text-xs font-semibold tracking-wide">
				{title}
			</p>
		</div>

		<div class="flex items-center gap-4">
			<form onsubmit={handleSearch} class="relative hidden md:block">
				<svg
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					class="pointer-events-none absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-slate-400"
				>
					<circle cx="11" cy="11" r="8"></circle>
					<path d="m21 21-4.3-4.3"></path>
				</svg>

				<input
					bind:value={search}
					type="search"
					placeholder="Search monitor..."
					class="h-10 w-52 rounded-md border-0 bg-white pr-4 pl-10 text-xs text-slate-700 ring-0 outline-none placeholder:text-slate-400 lg:w-64"
				/>
			</form>

			<div class="relative">
				<button
					type="button"
					onclick={() => {
						accountOpen = !accountOpen;
					}}
					class="flex h-10 w-10 items-center justify-center rounded-full border-2 border-white/20 bg-sky-500 text-sm font-bold text-white shadow"
				>
					{getInitial($authStore.user?.name)}
				</button>

				{#if accountOpen}
					<div
						class="absolute right-0 mt-3 w-56 overflow-hidden rounded-lg border border-slate-200 bg-white text-slate-700 shadow-xl"
					>
						<div class="border-b border-slate-100 px-4 py-3">
							<p class="truncate text-sm font-semibold text-slate-800">
								{$authStore.user?.name ?? 'User'}
							</p>

							<p class="mt-1 truncate text-xs text-slate-400">
								{$authStore.user?.email ?? ''}
							</p>
						</div>

						<button
							type="button"
							onclick={() => {
								accountOpen = false;
								void goto(resolve('/profile'));
							}}
							class="w-full px-4 py-3 text-left text-xs font-medium transition hover:bg-slate-50"
						>
							Profile
						</button>

						<button
							type="button"
							onclick={() => {
								void logout();
							}}
							class="w-full border-t border-slate-100 px-4 py-3 text-left text-xs font-medium text-red-500 transition hover:bg-red-50"
						>
							Logout
						</button>
					</div>
				{/if}
			</div>
		</div>
	</div>
</header>
