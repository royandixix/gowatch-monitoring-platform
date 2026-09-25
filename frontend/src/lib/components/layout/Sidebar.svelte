<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { authStore } from '$lib/stores/auth';

	let {
		open = false,
		onClose
	}: {
		open?: boolean;
		onClose: () => void;
	} = $props();

	type NavigationHref = '/dashboard' | '/monitors' | '/history' | '/profile';

	type NavigationItem = {
		name: string;
		href: NavigationHref;
		icon: 'dashboard' | 'monitor' | 'history' | 'profile';
	};

	const mainNavigation: NavigationItem[] = [
		{
			name: 'Dashboard',
			href: '/dashboard',
			icon: 'dashboard'
		},
		{
			name: 'Monitors',
			href: '/monitors',
			icon: 'monitor'
		},
		{
			name: 'History',
			href: '/history',
			icon: 'history'
		}
	];

	const accountNavigation: NavigationItem[] = [
		{
			name: 'Profile',
			href: '/profile',
			icon: 'profile'
		}
	];

	function isActive(href: string): boolean {
		if (href === '/dashboard') {
			return page.url.pathname === '/dashboard';
		}

		return page.url.pathname.startsWith(href);
	}

	async function navigate(href: NavigationHref): Promise<void> {
		onClose();
		await goto(resolve(href));
	}

	async function logout(): Promise<void> {
		authStore.logout();
		onClose();
		await goto(resolve('/login'));
	}
</script>

{#if open}
	<button
		type="button"
		aria-label="Tutup sidebar"
		class="fixed inset-0 z-40 bg-slate-950/40 backdrop-blur-[1px] lg:hidden"
		onclick={onClose}
	></button>
{/if}

<aside
	class={`fixed inset-y-0 left-0 z-50 flex w-64 flex-col border-r border-slate-200 bg-white shadow-sm transition-transform duration-200 lg:translate-x-0 ${
		open ? 'translate-x-0' : '-translate-x-full'
	}`}
>
	<div class="flex h-20 items-center px-6">
		<button
			type="button"
			class="flex items-center gap-3 text-left"
			onclick={() => {
				void navigate('/dashboard');
			}}
		>
			<div
				class="flex h-10 w-10 items-center justify-center rounded-lg bg-slate-800 text-sm font-bold text-white"
			>
				G
			</div>

			<div>
				<p class="text-sm font-bold tracking-wide text-slate-800">GOWATCH</p>

				<p class="mt-0.5 text-[10px] font-medium tracking-[0.14em] text-slate-400 uppercase">
					Monitoring Platform
				</p>
			</div>
		</button>
	</div>

	<div class="mx-5 border-t border-slate-200"></div>

	<nav class="flex-1 overflow-y-auto px-5 py-5">
		<p class="mb-3 text-[10px] font-bold tracking-[0.12em] text-slate-400 uppercase">Monitoring</p>

		<div class="space-y-1">
			{#each mainNavigation as item (item.href)}
				<button
					type="button"
					onclick={() => {
						void navigate(item.href);
					}}
					class={`group flex w-full items-center gap-3 rounded-md px-2 py-3 text-left text-xs font-semibold tracking-wide uppercase transition ${
						isActive(item.href)
							? 'text-sky-500'
							: 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
					}`}
				>
					<span
						class={`flex h-5 w-5 items-center justify-center ${
							isActive(item.href) ? 'text-sky-500' : 'text-slate-400'
						}`}
					>
						{#if item.icon === 'dashboard'}
							<svg
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								class="h-4 w-4"
							>
								<rect x="3" y="3" width="7" height="7" rx="1"></rect>
								<rect x="14" y="3" width="7" height="7" rx="1"></rect>
								<rect x="3" y="14" width="7" height="7" rx="1"></rect>
								<rect x="14" y="14" width="7" height="7" rx="1"></rect>
							</svg>
						{:else if item.icon === 'monitor'}
							<svg
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								class="h-4 w-4"
							>
								<rect x="3" y="4" width="18" height="13" rx="2"></rect>
								<path d="M8 21h8"></path>
								<path d="M12 17v4"></path>
							</svg>
						{:else if item.icon === 'history'}
							<svg
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								class="h-4 w-4"
							>
								<path d="M3 12a9 9 0 1 0 3-6.7"></path>
								<path d="M3 4v6h6"></path>
								<path d="M12 7v5l3 2"></path>
							</svg>
						{/if}
					</span>

					{item.name}
				</button>
			{/each}
		</div>

		<div class="my-6 border-t border-slate-200"></div>

		<p class="mb-3 text-[10px] font-bold tracking-[0.12em] text-slate-400 uppercase">Account</p>

		<div class="space-y-1">
			{#each accountNavigation as item (item.href)}
				<button
					type="button"
					onclick={() => {
						void navigate(item.href);
					}}
					class={`group flex w-full items-center gap-3 rounded-md px-2 py-3 text-left text-xs font-semibold tracking-wide uppercase transition ${
						isActive(item.href)
							? 'text-sky-500'
							: 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
					}`}
				>
					<span
						class={`flex h-5 w-5 items-center justify-center ${
							isActive(item.href) ? 'text-sky-500' : 'text-slate-400'
						}`}
					>
						<svg
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							class="h-4 w-4"
						>
							<circle cx="12" cy="8" r="4"></circle>
							<path d="M4 21a8 8 0 0 1 16 0"></path>
						</svg>
					</span>

					{item.name}
				</button>
			{/each}

			<button
				type="button"
				onclick={() => {
					void logout();
				}}
				class="flex w-full items-center gap-3 rounded-md px-2 py-3 text-left text-xs font-semibold tracking-wide text-slate-600 uppercase transition hover:bg-red-50 hover:text-red-500"
			>
				<span class="flex h-5 w-5 items-center justify-center">
					<svg
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						class="h-4 w-4"
					>
						<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
						<path d="m16 17 5-5-5-5"></path>
						<path d="M21 12H9"></path>
					</svg>
				</span>

				Logout
			</button>
		</div>
	</nav>

	<div class="border-t border-slate-200 px-6 py-5">
		<div class="flex items-center gap-2">
			<span class="h-2 w-2 rounded-full bg-emerald-500"></span>

			<div>
				<p class="text-xs font-semibold text-slate-700">GoWatch API</p>

				<p class="text-[10px] text-slate-400">Monitoring service</p>
			</div>
		</div>
	</div>
</aside>
