<script lang="ts">
	import { onMount } from 'svelte';

	import { goto } from '$app/navigation';

	import { resolve } from '$app/paths';

	import { page } from '$app/state';

	import { getMonitorResults, getMonitorStats } from '$lib/api/history';

	import { getMonitor, updateMonitor } from '$lib/api/monitors';

	import PageHeader from '$lib/components/layout/PageHeader.svelte';

	import MonitorForm from '$lib/components/monitor/MonitorForm.svelte';

	import MonitorStatusBadge from '$lib/components/monitor/MonitorStatusBadge.svelte';

	import StatCard from '$lib/components/ui/StatCard.svelte';

	import type { MonitorResultRecord, MonitorStats } from '$lib/types/history';

	import type { CreateMonitorPayload, Monitor } from '$lib/types/monitor';

	import { showError, showSuccess } from '$lib/utils/alert';

	import { formatDate, formatInterval } from '$lib/utils/format';

	const monitorId = Number(page.params.id);

	let monitor = $state<Monitor | null>(null);

	let stats = $state<MonitorStats | null>(null);

	let results = $state<MonitorResultRecord[]>([]);

	let loading = $state(true);

	let refreshing = $state(false);

	let saving = $state(false);

	let showEditForm = $state(false);

	let errorMessage = $state('');

	const currentStatus = $derived.by(() => {
		if (!monitor) {
			return 'WAITING';
		}

		if (!monitor.is_active) {
			return 'PAUSED';
		}

		return stats?.last_status || 'WAITING';
	});

	const chartResults = $derived.by(() => {
		return [...results.slice(0, 30)].reverse();
	});

	const recentResults = $derived(results.slice(0, 20));

	const chartMaxResponse = $derived.by(() => {
		if (chartResults.length === 0) {
			return 1;
		}

		return Math.max(...chartResults.map((result) => result.response_time_ms), 1);
	});

	const chartPoints = $derived.by(() => {
		if (chartResults.length === 0) {
			return '';
		}

		return chartResults
			.map((result, index) => {
				const x = getChartX(index, chartResults.length);

				const y = getChartY(result.response_time_ms, chartMaxResponse);

				return `${x},${y}`;
			})
			.join(' ');
	});

	onMount(() => {
		void loadDetail(true);
	});

	function getErrorMessage(error: unknown, fallback: string): string {
		if (error instanceof Error) {
			return error.message;
		}

		return fallback;
	}

	function getChartX(index: number, count: number): number {
		if (count <= 1) {
			return 500;
		}

		return 30 + (index / (count - 1)) * 940;
	}

	function getChartY(value: number, maxValue: number): number {
		const normalized = Math.max(value, 0) / Math.max(maxValue, 1);

		return 190 - normalized * 150;
	}

	function getStatusTone(value: string): 'red' | 'orange' | 'emerald' {
		if (value === 'DOWN') {
			return 'red';
		}

		if (value === 'UP') {
			return 'emerald';
		}

		return 'orange';
	}

	function getStatusIcon(value: string): 'up' | 'down' | 'performance' {
		if (value === 'DOWN') {
			return 'down';
		}

		if (value === 'UP') {
			return 'up';
		}

		return 'performance';
	}

	function formatResponseTime(value?: number): string {
		if (!value || value <= 0) {
			return '-';
		}

		return `${value} ms`;
	}

	async function loadDetail(initial = false): Promise<void> {
		if (!Number.isInteger(monitorId) || monitorId <= 0) {
			errorMessage = 'ID monitor tidak valid.';

			loading = false;

			return;
		}

		if (initial) {
			loading = true;
		} else {
			refreshing = true;
		}

		errorMessage = '';

		try {
			const [monitorResponse, historyResponse, statsResponse] = await Promise.all([
				getMonitor(monitorId),
				getMonitorResults(monitorId, 100),
				getMonitorStats(monitorId)
			]);

			monitor = monitorResponse.monitor;

			results = historyResponse.results;

			stats = statsResponse.stats;
		} catch (error) {
			errorMessage = getErrorMessage(error, 'Gagal memuat detail monitor.');

			await showError('Gagal memuat monitor', errorMessage);
		} finally {
			loading = false;
			refreshing = false;
		}
	}

	async function refreshDetail(): Promise<void> {
		await loadDetail(false);
	}

	function openEdit(): void {
		showEditForm = true;

		window.scrollTo({
			top: 0,
			behavior: 'smooth'
		});
	}

	function closeEdit(): void {
		showEditForm = false;
	}

	async function saveMonitor(payload: CreateMonitorPayload): Promise<void> {
		if (!monitor || saving) {
			return;
		}

		saving = true;

		try {
			const response = await updateMonitor(monitor.id, payload);

			monitor = response.monitor;

			showEditForm = false;

			await showSuccess('Monitor diperbarui', 'Konfigurasi monitor berhasil disimpan.');
		} catch (error) {
			await showError(
				'Gagal memperbarui monitor',
				getErrorMessage(error, 'Terjadi kesalahan saat memperbarui monitor.')
			);
		} finally {
			saving = false;
		}
	}

	async function backToMonitors(): Promise<void> {
		await goto(resolve('/monitors'));
	}
</script>

<svelte:head>
	<title>
		{monitor?.name ?? 'Monitor Detail'} | GoWatch
	</title>
</svelte:head>

<PageHeader
	title={monitor?.name ?? 'Monitor Detail'}
	description={monitor?.url ?? 'Informasi detail, statistik, dan riwayat monitoring.'}
/>

<div class="-mt-14 px-4 pb-10 sm:px-6 lg:px-8">
	<div class="mx-auto max-w-7xl">
		{#if errorMessage && !loading}
			<div class="mb-6 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-600">
				{errorMessage}
			</div>
		{/if}

		{#if loading}
			<div class="flex min-h-72 items-center justify-center rounded-lg bg-white shadow-lg">
				<div class="flex items-center gap-3">
					<div
						class="h-5 w-5 animate-spin rounded-full border-2 border-slate-200 border-t-sky-500"
					></div>

					<p class="text-sm text-slate-400">Memuat detail monitor...</p>
				</div>
			</div>
		{:else if monitor}
			<div class="mb-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end">
				<button
					type="button"
					onclick={() => {
						void refreshDetail();
					}}
					disabled={refreshing}
					class="inline-flex h-10 items-center justify-center gap-2 rounded-lg border border-slate-200 bg-white px-4 text-xs font-semibold text-slate-600 shadow-sm transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50"
				>
					<svg
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.8"
						stroke-linecap="round"
						stroke-linejoin="round"
						class={`h-4 w-4 ${refreshing ? 'animate-spin' : ''}`}
						aria-hidden="true"
					>
						<path d="M20 11a8.1 8.1 0 0 0-15.5-2M4 4v5h5"></path>

						<path d="M4 13a8.1 8.1 0 0 0 15.5 2M20 20v-5h-5"></path>
					</svg>

					Refresh
				</button>

				<button
					type="button"
					onclick={openEdit}
					class="inline-flex h-10 items-center justify-center gap-2 rounded-lg bg-sky-500 px-4 text-xs font-semibold text-white shadow-sm transition hover:bg-sky-600"
				>
					<svg
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.8"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="h-4 w-4"
						aria-hidden="true"
					>
						<path d="M12 20h9"></path>

						<path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L8 18l-4 1 1-4Z"></path>
					</svg>

					Edit Monitor
				</button>
			</div>

			{#if showEditForm}
				<div class="mb-6">
					{#key monitor.updated_at}
						<MonitorForm {monitor} loading={saving} onSubmit={saveMonitor} onCancel={closeEdit} />
					{/key}
				</div>
			{/if}

			<div class="grid gap-5 sm:grid-cols-2 xl:grid-cols-4">
				<StatCard
					title="Current Status"
					value={currentStatus}
					subtitle={monitor.is_active
						? 'Status pengecekan terakhir'
						: 'Monitoring sedang dihentikan'}
					tone={getStatusTone(currentStatus)}
					icon={getStatusIcon(currentStatus)}
				/>

				<StatCard
					title="Uptime"
					value={`${(stats?.uptime_percentage ?? 0).toFixed(2)}%`}
					subtitle={`${stats?.total_up ?? 0} UP dari ${stats?.total_checks ?? 0} checks`}
					tone="emerald"
					icon="up"
				/>

				<StatCard
					title="Average Response"
					value={formatResponseTime(stats?.average_response_time_ms)}
					subtitle="Rata-rata seluruh pengecekan"
					tone="sky"
					icon="performance"
				/>

				<StatCard
					title="Total Checks"
					value={stats?.total_checks ?? 0}
					subtitle="Jumlah pengecekan tersimpan"
					tone="pink"
					icon="monitor"
				/>
			</div>

			<div class="mt-6 grid gap-6 xl:grid-cols-[360px_minmax(0,1fr)]">
				<div class="rounded-lg bg-white p-6 shadow-lg">
					<div class="flex items-start justify-between gap-4 border-b border-slate-100 pb-5">
						<div>
							<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">
								Configuration
							</p>

							<h2 class="mt-1 text-base font-semibold text-slate-800">Monitor Information</h2>
						</div>

						<MonitorStatusBadge status={currentStatus} active={monitor.is_active} />
					</div>

					<div class="mt-6 space-y-5">
						<div>
							<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">Name</p>

							<p class="mt-1 text-sm font-semibold text-slate-700">
								{monitor.name}
							</p>
						</div>

						<div>
							<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">URL</p>

							<a
								href={monitor.url}
								target="_blank"
								rel="external noreferrer"
								class="mt-1 block text-sm font-medium break-all text-sky-500 transition hover:text-sky-600 hover:underline"
							>
								{monitor.url}
							</a>
						</div>

						<div class="grid grid-cols-2 gap-4">
							<div>
								<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">Type</p>

								<p class="mt-1 text-sm font-semibold text-slate-700 uppercase">
									{monitor.monitor_type}
								</p>
							</div>

							<div>
								<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">Interval</p>

								<p class="mt-1 text-sm font-semibold text-slate-700">
									{formatInterval(monitor.interval_seconds)}
								</p>
							</div>
						</div>

						<div class="grid grid-cols-2 gap-4">
							<div>
								<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">
									Last Check
								</p>

								<p class="mt-1 text-xs font-medium text-slate-600">
									{stats?.last_checked_at ? formatDate(stats.last_checked_at) : '-'}
								</p>
							</div>

							<div>
								<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">Created</p>

								<p class="mt-1 text-xs font-medium text-slate-600">
									{formatDate(monitor.created_at)}
								</p>
							</div>
						</div>
					</div>

					<button
						type="button"
						onclick={() => {
							void backToMonitors();
						}}
						class="mt-7 flex h-10 w-full items-center justify-center gap-2 rounded-lg bg-slate-800 text-xs font-semibold text-white transition hover:bg-slate-700"
					>
						<svg
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="h-4 w-4"
							aria-hidden="true"
						>
							<path d="m15 18-6-6 6-6"></path>
						</svg>

						Kembali ke Monitors
					</button>
				</div>

				<div class="overflow-hidden rounded-lg bg-slate-800 shadow-lg">
					<div class="border-b border-slate-700 px-6 py-5">
						<div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
							<div>
								<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">
									Performance
								</p>

								<h2 class="mt-1 text-base font-semibold text-white">Response Time</h2>

								<p class="mt-1 text-xs text-slate-400">30 pengecekan terbaru</p>
							</div>

							<div class="grid grid-cols-2 gap-x-6 gap-y-3 sm:grid-cols-4">
								<div>
									<p class="text-[9px] font-bold tracking-wide text-slate-500 uppercase">Last</p>

									<p class="mt-1 text-xs font-semibold text-white">
										{formatResponseTime(stats?.last_response_time_ms)}
									</p>
								</div>

								<div>
									<p class="text-[9px] font-bold tracking-wide text-slate-500 uppercase">HTTP</p>

									<p class="mt-1 text-xs font-semibold text-white">
										{stats?.last_http_status ? stats.last_http_status : '-'}
									</p>
								</div>

								<div>
									<p class="text-[9px] font-bold tracking-wide text-slate-500 uppercase">
										Total UP
									</p>

									<p class="mt-1 text-xs font-semibold text-emerald-400">
										{stats?.total_up ?? 0}
									</p>
								</div>

								<div>
									<p class="text-[9px] font-bold tracking-wide text-slate-500 uppercase">
										Total DOWN
									</p>

									<p class="mt-1 text-xs font-semibold text-red-400">
										{stats?.total_down ?? 0}
									</p>
								</div>
							</div>
						</div>
					</div>

					<div class="p-6">
						{#if chartResults.length === 0}
							<div
								class="flex h-64 items-center justify-center rounded-lg border border-dashed border-slate-600"
							>
								<p class="text-sm text-slate-400">Belum ada data untuk grafik.</p>
							</div>
						{:else}
							<div class="mb-4 flex justify-between text-[10px] font-medium text-slate-500">
								<span>
									{chartMaxResponse} ms
								</span>

								<span>
									Latest response:
									{formatResponseTime(stats?.last_response_time_ms)}
								</span>
							</div>

							<div class="relative h-64 w-full">
								<svg
									viewBox="0 0 1000 220"
									preserveAspectRatio="none"
									class="h-full w-full"
									role="img"
									aria-label="Grafik response time monitor"
								>
									<line
										x1="30"
										y1="40"
										x2="970"
										y2="40"
										stroke="currentColor"
										stroke-opacity="0.08"
										class="text-white"
									></line>

									<line
										x1="30"
										y1="90"
										x2="970"
										y2="90"
										stroke="currentColor"
										stroke-opacity="0.08"
										class="text-white"
									></line>

									<line
										x1="30"
										y1="140"
										x2="970"
										y2="140"
										stroke="currentColor"
										stroke-opacity="0.08"
										class="text-white"
									></line>

									<line
										x1="30"
										y1="190"
										x2="970"
										y2="190"
										stroke="currentColor"
										stroke-opacity="0.08"
										class="text-white"
									></line>

									<polyline
										points={chartPoints}
										fill="none"
										stroke="currentColor"
										stroke-width="3"
										stroke-linecap="round"
										stroke-linejoin="round"
										vector-effect="non-scaling-stroke"
										class="text-sky-400"
									></polyline>

									{#each chartResults as result, index (result.id)}
										<circle
											cx={getChartX(index, chartResults.length)}
											cy={getChartY(result.response_time_ms, chartMaxResponse)}
											r="4"
											fill="currentColor"
											vector-effect="non-scaling-stroke"
											class={result.status === 'DOWN' ? 'text-red-400' : 'text-sky-300'}
										>
											<title>
												{result.status} -
												{result.response_time_ms} ms
											</title>
										</circle>
									{/each}
								</svg>
							</div>

							<div class="mt-3 flex items-center justify-between text-[10px] text-slate-500">
								<span> Lebih lama </span>

								<span> Lebih baru </span>
							</div>
						{/if}
					</div>
				</div>
			</div>

			<div class="mt-6 overflow-hidden rounded-lg bg-white shadow-lg">
				<div
					class="flex flex-col gap-2 border-b border-slate-100 px-6 py-5 sm:flex-row sm:items-center sm:justify-between"
				>
					<div>
						<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">History</p>

						<h2 class="mt-1 text-base font-semibold text-slate-800">Recent Checks</h2>
					</div>

					<p class="text-xs text-slate-400">20 pengecekan terbaru</p>
				</div>

				{#if recentResults.length === 0}
					<div class="flex min-h-56 items-center justify-center p-8">
						<p class="text-sm text-slate-400">Belum ada hasil pengecekan.</p>
					</div>
				{:else}
					<div class="overflow-x-auto">
						<table class="w-full min-w-[900px]">
							<thead class="bg-slate-50">
								<tr>
									<th
										class="px-6 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
									>
										Status
									</th>

									<th
										class="px-6 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
									>
										HTTP
									</th>

									<th
										class="px-6 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
									>
										Response
									</th>

									<th
										class="px-6 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
									>
										Checked At
									</th>

									<th
										class="px-6 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
									>
										Error
									</th>
								</tr>
							</thead>

							<tbody class="divide-y divide-slate-100">
								{#each recentResults as result (result.id)}
									<tr class="transition hover:bg-slate-50">
										<td class="px-6 py-4">
											<MonitorStatusBadge status={result.status} />
										</td>

										<td class="px-6 py-4">
											<span
												class={`text-xs font-semibold ${
													result.http_status >= 200 && result.http_status < 300
														? 'text-emerald-600'
														: result.http_status >= 400
															? 'text-red-500'
															: 'text-slate-500'
												}`}
											>
												{result.http_status || '-'}
											</span>
										</td>

										<td class="px-6 py-4 text-xs font-medium text-slate-600">
											{formatResponseTime(result.response_time_ms)}
										</td>

										<td class="px-6 py-4 text-xs text-slate-500">
											{formatDate(result.checked_at)}
										</td>

										<td class="px-6 py-4">
											<p
												class={`max-w-[320px] truncate text-xs ${
													result.error ? 'text-red-500' : 'text-slate-300'
												}`}
												title={result.error || undefined}
											>
												{result.error || '-'}
											</p>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>
