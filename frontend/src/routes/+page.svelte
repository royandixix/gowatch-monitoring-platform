<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { getToken } from '$lib/utils/storage';

	onMount(() => {
		const destination: '/dashboard' | '/login' = getToken() ? '/dashboard' : '/login';

		void goto(resolve(destination), {
			replaceState: true
		});
	});
</script>

<svelte:head>
	<title>GoWatch</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-slate-50">
	<div class="flex items-center gap-3 text-slate-500">
		<div
			class="h-5 w-5 animate-spin rounded-full border-2 border-slate-300 border-t-blue-600"
		></div>

		<span class="text-sm font-medium"> Memuat GoWatch... </span>
	</div>
</div>
