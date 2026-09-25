<script lang="ts">
	let {
		status,
		active
	}: {
		status?: string;
		active?: boolean;
	} = $props();

	const label = $derived.by(() => {
		if (status) {
			return status.toUpperCase();
		}

		return active ? 'ACTIVE' : 'PAUSED';
	});

	const badgeClass = $derived.by(() => {
		switch (label) {
			case 'UP':
			case 'ACTIVE':
				return 'bg-emerald-50 text-emerald-600 ring-emerald-200';

			case 'DOWN':
				return 'bg-red-50 text-red-600 ring-red-200';

			case 'PAUSED':
			case 'INACTIVE':
				return 'bg-slate-100 text-slate-500 ring-slate-200';

			default:
				return 'bg-amber-50 text-amber-600 ring-amber-200';
		}
	});

	const dotClass = $derived.by(() => {
		switch (label) {
			case 'UP':
			case 'ACTIVE':
				return 'bg-emerald-500';

			case 'DOWN':
				return 'bg-red-500';

			case 'PAUSED':
			case 'INACTIVE':
				return 'bg-slate-400';

			default:
				return 'bg-amber-500';
		}
	});
</script>

<span
	class={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[10px] font-bold ring-1 ring-inset ${badgeClass}`}
>
	<span class={`h-1.5 w-1.5 rounded-full ${dotClass}`}></span>

	{label}
</span>
