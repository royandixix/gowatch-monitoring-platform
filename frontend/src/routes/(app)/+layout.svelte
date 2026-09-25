<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { get } from 'svelte/store';

	import Header from '$lib/components/layout/Header.svelte';
	import Sidebar from '$lib/components/layout/Sidebar.svelte';

	import { authStore } from '$lib/stores/auth';

	let { children } = $props();

	let ready = $state(false);

	let mobileSidebarOpen = $state(false);

	onMount(() => {
		void initialize();
	});

	async function initialize(): Promise<void> {
		await authStore.initialize();

		const state = get(authStore);

		if (!state.isAuthenticated) {
			await goto(resolve('/login'), {
				replaceState: true
			});

			return;
		}

		ready = true;
	}
</script>

{#if ready}
	<div class="min-h-screen bg-slate-100">
		<Sidebar
			open={mobileSidebarOpen}
			onClose={() => {
				mobileSidebarOpen = false;
			}}
		/>

		<div class="min-h-screen lg:pl-64">
			<Header
				onMenuClick={() => {
					mobileSidebarOpen = true;
				}}
			/>

			<main class="min-h-[calc(100vh-4rem)]">
				{@render children()}
			</main>
		</div>
	</div>
{:else}
	<div class="flex min-h-screen items-center justify-center bg-slate-100">
		<div class="flex items-center gap-3">
			<div
				class="h-5 w-5 animate-spin rounded-full border-2 border-slate-300 border-t-sky-500"
			></div>

			<p class="text-sm font-medium text-slate-500">Memeriksa sesi...</p>
		</div>
	</div>
{/if}
