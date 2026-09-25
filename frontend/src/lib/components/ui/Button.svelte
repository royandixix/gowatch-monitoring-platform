<script lang="ts">
	type Variant = 'primary' | 'secondary' | 'danger' | 'ghost';

	let {
		type = 'button',
		variant = 'primary',
		disabled = false,
		loading = false,
		onclick,
		children
	}: {
		type?: 'button' | 'submit' | 'reset';
		variant?: Variant;
		disabled?: boolean;
		loading?: boolean;
		onclick?: (event: MouseEvent) => void;
		children: import('svelte').Snippet;
	} = $props();

	const variantClass = $derived.by(() => {
		switch (variant) {
			case 'secondary':
				return 'border border-slate-200 bg-white text-slate-700 hover:bg-slate-50';

			case 'danger':
				return 'bg-red-500 text-white hover:bg-red-600';

			case 'ghost':
				return 'bg-transparent text-slate-600 hover:bg-slate-100';

			default:
				return 'bg-sky-500 text-white hover:bg-sky-600';
		}
	});
</script>

<button
	{type}
	{onclick}
	disabled={disabled || loading}
	class={`inline-flex h-10 items-center justify-center gap-2 rounded-md px-4 text-xs font-semibold transition disabled:cursor-not-allowed disabled:opacity-60 ${variantClass}`}
>
	{#if loading}
		<span class="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"
		></span>
	{/if}

	{@render children()}
</button>
