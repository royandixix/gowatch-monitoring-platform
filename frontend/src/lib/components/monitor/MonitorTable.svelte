<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	import MonitorStatusBadge from '$lib/components/monitor/MonitorStatusBadge.svelte';

	import type { Monitor } from '$lib/types/monitor';

	import { formatDate, formatInterval } from '$lib/utils/format';

	let {
		monitors,
		onEdit,
		onDelete
	}: {
		monitors: Monitor[];
		onEdit: (monitor: Monitor) => void;
		onDelete: (monitor: Monitor) => void;
	} = $props();

	async function openDetail(monitor: Monitor): Promise<void> {
		await goto(resolve(`/monitors/${monitor.id}`));
	}
</script>

<div class="overflow-hidden rounded-md bg-white shadow-lg">
	<div class="overflow-x-auto">
		<table class="w-full min-w-[820px]">
			<thead class="bg-slate-50">
				<tr>
					<th
						class="px-6 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
					>
						Monitor
					</th>

					<th
						class="px-6 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
					>
						Tipe
					</th>

					<th
						class="px-6 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
					>
						Interval
					</th>

					<th
						class="px-6 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
					>
						Status
					</th>

					<th
						class="px-6 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
					>
						Dibuat
					</th>

					<th
						class="px-6 py-4 text-right text-[10px] font-bold tracking-wide text-slate-400 uppercase"
					>
						Aksi
					</th>
				</tr>
			</thead>

			<tbody class="divide-y divide-slate-100">
				{#each monitors as monitor (monitor.id)}
					<tr class="transition-colors duration-150 hover:bg-slate-50">
						<td class="px-6 py-4">
							<div class="min-w-0">
								<p class="text-sm font-semibold text-slate-700">
									{monitor.name}
								</p>

								<p class="mt-1 max-w-xs truncate text-xs text-slate-400">
									{monitor.url}
								</p>
							</div>
						</td>

						<td class="px-6 py-4">
							<span class="text-xs font-medium text-slate-500 uppercase">
								{monitor.monitor_type}
							</span>
						</td>

						<td class="px-6 py-4 text-xs text-slate-500">
							{formatInterval(monitor.interval_seconds)}
						</td>

						<td class="px-6 py-4">
							<MonitorStatusBadge active={monitor.is_active} />
						</td>

						<td class="px-6 py-4 text-xs text-slate-500">
							{formatDate(monitor.created_at)}
						</td>

						<td class="px-6 py-4">
							<div class="flex items-center justify-end gap-2">
								<div class="group relative">
									<button
										type="button"
										aria-label={`Lihat detail ${monitor.name}`}
										onclick={() => {
											void openDetail(monitor);
										}}
										class="flex h-9 w-9 items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-500 shadow-sm transition duration-200 hover:border-slate-300 hover:bg-slate-100 hover:text-slate-800"
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
											<path
												d="M2.25 12s3.75-6.75 9.75-6.75S21.75 12 21.75 12 18 18.75 12 18.75 2.25 12 2.25 12Z"
											></path>

											<circle cx="12" cy="12" r="3"></circle>
										</svg>
									</button>

									<div
										class="pointer-events-none absolute bottom-full left-1/2 z-20 mb-2 -translate-x-1/2 rounded bg-slate-800 px-2 py-1 text-[10px] font-medium whitespace-nowrap text-white opacity-0 shadow-lg transition-opacity group-hover:opacity-100"
									>
										Detail
									</div>
								</div>

								<div class="group relative">
									<button
										type="button"
										aria-label={`Edit ${monitor.name}`}
										onclick={() => {
											onEdit(monitor);
										}}
										class="flex h-9 w-9 items-center justify-center rounded-lg border border-sky-100 bg-sky-50 text-sky-500 shadow-sm transition duration-200 hover:border-sky-200 hover:bg-sky-100 hover:text-sky-700"
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

									<div
										class="pointer-events-none absolute bottom-full left-1/2 z-20 mb-2 -translate-x-1/2 rounded bg-slate-800 px-2 py-1 text-[10px] font-medium whitespace-nowrap text-white opacity-0 shadow-lg transition-opacity group-hover:opacity-100"
									>
										Edit
									</div>
								</div>

								<div class="group relative">
									<button
										type="button"
										aria-label={`Hapus ${monitor.name}`}
										onclick={() => {
											onDelete(monitor);
										}}
										class="flex h-9 w-9 items-center justify-center rounded-lg border border-red-100 bg-red-50 text-red-500 shadow-sm transition duration-200 hover:border-red-200 hover:bg-red-100 hover:text-red-600"
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

									<div
										class="pointer-events-none absolute bottom-full left-1/2 z-20 mb-2 -translate-x-1/2 rounded bg-slate-800 px-2 py-1 text-[10px] font-medium whitespace-nowrap text-white opacity-0 shadow-lg transition-opacity group-hover:opacity-100"
									>
										Hapus
									</div>
								</div>
							</div>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</div>
