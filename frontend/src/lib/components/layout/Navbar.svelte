<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	import { authStore } from '$lib/stores/auth';

	let {
		onMenuClick
	}: {
		onMenuClick: () => void;
	} = $props();

	function getInitial(name: string | undefined): string {
		if (!name) {
			return 'U';
		}

		return name.trim().charAt(0).toUpperCase();
	}

	async function handleLogout(): Promise<void> {
		authStore.logout();

		await goto(resolve('/login'));
	}
</script>

<header
	class="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-slate-200 bg-white/95 px-4 backdrop-blur sm:px-6 lg:px-8"
>
	<div class="flex items-center gap-3">
		<button
			type="button"
			aria-label="Buka sidebar"
			class="flex h-10 w-10 items-center justify-center rounded-lg border border-slate-200 text-slate-600 hover:bg-slate-50 lg:hidden"
			onclick={onMenuClick}
		>
			<span class="text-xl"> ☰ </span>
		</button>

		<div>
			<p class="text-sm font-semibold text-slate-950">GoWatch Dashboard</p>

			<p class="hidden text-xs text-slate-500 sm:block">Website & API Monitoring</p>
		</div>
	</div>

	<div class="flex items-center gap-3">
		<div class="hidden text-right sm:block">
			<p class="text-sm font-semibold text-slate-800">
				{$authStore.user?.name ?? 'User'}
			</p>

			<p class="text-xs text-slate-500">
				{$authStore.user?.email ?? ''}
			</p>
		</div>

		<div
			class="flex h-10 w-10 items-center justify-center rounded-full bg-slate-900 text-sm font-bold text-white"
		>
			{getInitial($authStore.user?.name)}
		</div>

		<button
			type="button"
			onclick={handleLogout}
			class="hidden rounded-lg border border-slate-200 px-3 py-2 text-xs font-semibold text-slate-600 transition hover:bg-slate-50 sm:block"
		>
			Logout
		</button>
	</div>
</header>
