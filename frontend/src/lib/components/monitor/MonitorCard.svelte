<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	import MonitorStatusBadge from '$lib/components/monitor/MonitorStatusBadge.svelte';

	import type { Monitor } from '$lib/types/monitor';

	import { formatInterval } from '$lib/utils/format';

	let {
		monitor,
		onEdit,
		onDelete
	}: {
		monitor: Monitor;
		onEdit: (monitor: Monitor) => void;
		onDelete: (monitor: Monitor) => void;
	} = $props();

	async function openDetail(): Promise<void> {
		await goto(resolve(`/monitors/${monitor.id}`));
	}
</script>

<div class="rounded-lg border border-slate-100 bg-white p-5 shadow-sm">
	<div class="flex items-start justify-between gap-4">
		<div class="min-w-0">
			<p class="truncate text-sm font-semibold text-slate-700">
				{monitor.name}
			</p>

			<p class="mt-1 truncate text-xs text-slate-400">
				{monitor.url}
			</p>
		</div>

		<MonitorStatusBadge active={monitor.is_active} />
	</div>

	<div class="mt-5 grid grid-cols-2 gap-3 border-t border-slate-100 pt-4">
		<div>
			<p class="text-[10px] font-semibold tracking-wide text-slate-400 uppercase">Tipe</p>

			<p class="mt-1 text-xs font-semibold text-slate-600 uppercase">
				{monitor.monitor_type}
			</p>
		</div>

		<div>
			<p class="text-[10px] font-semibold tracking-wide text-slate-400 uppercase">Interval</p>

			<p class="mt-1 text-xs font-semibold text-slate-600">
				{formatInterval(monitor.interval_seconds)}
			</p>
		</div>
	</div>

	<div class="mt-5 flex items-center justify-end gap-2 border-t border-slate-100 pt-4">
		<button
			type="button"
			aria-label={`Lihat detail ${monitor.name}`}
			onclick={() => {
				void openDetail();
			}}
			class="flex h-9 w-9 items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-500 transition hover:bg-slate-100 hover:text-slate-800"
		>
			<svg
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="1.8"
				stroke-linecap="round"
				stroke-linejoin="round"
				class="h-[18px] w-[18px]"
				aria-hidden="true"
			>
				<path d="M2.25 12s3.75-6.75 9.75-6.75S21.75 12 21.75 12 18 18.75 12 18.75 2.25 12 2.25 12Z"
				></path>

				<circle cx="12" cy="12" r="3"></circle>
			</svg>
		</button>

		<button
			type="button"
			aria-label={`Edit ${monitor.name}`}
			onclick={() => {
				onEdit(monitor);
			}}
			class="flex h-9 w-9 items-center justify-center rounded-lg border border-sky-100 bg-sky-50 text-sky-500 transition hover:bg-sky-100 hover:text-sky-700"
		>
			<svg
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="1.8"
				stroke-linecap="round"
				stroke-linejoin="round"
				class="h-[18px] w-[18px]"
				aria-hidden="true"
			>
				<path d="M12 20h9"></path>

				<path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L8 18l-4 1 1-4Z"></path>
			</svg>
		</button>

		<button
			type="button"
			aria-label={`Hapus ${monitor.name}`}
			onclick={() => {
				onDelete(monitor);
			}}
			class="flex h-9 w-9 items-center justify-center rounded-lg border border-red-100 bg-red-50 text-red-500 transition hover:bg-red-100 hover:text-red-600"
		>
			<svg
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="1.8"
				stroke-linecap="round"
				stroke-linejoin="round"
				class="h-[18px] w-[18px]"
				aria-hidden="true"
			>
				<path d="M3 6h18"></path>

				<path d="M8 6V4h8v2"></path>

				<path d="M19 6l-1 14H6L5 6"></path>

				<path d="M10 11v5"></path>

				<path d="M14 11v5"></path>
			</svg>
		</button>
	</div>
</div>
