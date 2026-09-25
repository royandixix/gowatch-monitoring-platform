<script lang="ts">
	type Tone = 'red' | 'orange' | 'pink' | 'sky' | 'emerald';

	let {
		title,
		value,
		subtitle,
		tone = 'sky',
		icon = 'monitor'
	}: {
		title: string;
		value: string | number;
		subtitle?: string;
		tone?: Tone;
		icon?: 'monitor' | 'up' | 'down' | 'performance';
	} = $props();

	const iconClass = $derived.by(() => {
		switch (tone) {
			case 'red':
				return 'bg-red-500';

			case 'orange':
				return 'bg-orange-500';

			case 'pink':
				return 'bg-pink-500';

			case 'emerald':
				return 'bg-emerald-500';

			default:
				return 'bg-sky-500';
		}
	});
</script>

<div class="rounded-md bg-white px-5 py-4 shadow-lg">
	<div class="flex items-center justify-between gap-5">
		<div>
			<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">
				{title}
			</p>

			<p class="mt-1 text-xl font-semibold text-slate-700">
				{value}
			</p>
		</div>

		<div
			class={`flex h-11 w-11 shrink-0 items-center justify-center rounded-full text-white shadow ${iconClass}`}
		>
			{#if icon === 'monitor'}
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="h-5 w-5">
					<rect x="3" y="4" width="18" height="13" rx="2"></rect>
					<path d="M8 21h8"></path>
					<path d="M12 17v4"></path>
				</svg>
			{:else if icon === 'up'}
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="h-5 w-5">
					<path d="m18 15-6-6-6 6"></path>
				</svg>
			{:else if icon === 'down'}
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="h-5 w-5">
					<path d="m6 9 6 6 6-6"></path>
				</svg>
			{:else}
				<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="h-5 w-5">
					<circle cx="12" cy="12" r="9"></circle>
					<path d="M12 7v5l3 2"></path>
				</svg>
			{/if}
		</div>
	</div>

	{#if subtitle}
		<p class="mt-3 text-[11px] text-slate-400">
			{subtitle}
		</p>
	{/if}
</div>
