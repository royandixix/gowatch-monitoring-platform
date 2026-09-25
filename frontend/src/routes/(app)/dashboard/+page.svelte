<script lang="ts">
	import { onMount } from 'svelte';

	import { goto } from '$app/navigation';

	import { resolve } from '$app/paths';

	import { getDashboard } from '$lib/api/dashboard';

	import PageHeader from '$lib/components/layout/PageHeader.svelte';

	import StatCard from '$lib/components/ui/StatCard.svelte';

	import type { DashboardResponse } from '$lib/types/dashboard';

	import { showError } from '$lib/utils/alert';

	import { formatDate } from '$lib/utils/format';

	let dashboard = $state<DashboardResponse | null>(null);

	let loading = $state(true);

	let refreshing = $state(false);

	const summary = $derived(
		dashboard?.summary ?? {
			total_monitors: 0,
			up_monitors: 0,
			down_monitors: 0,
			paused_monitors: 0,
			waiting_monitors: 0,
			total_checks_24h: 0,
			total_checks_7d: 0,
			uptime_24h_percentage: 0,
			uptime_7d_percentage: 0,
			average_response_time_24h_ms: 0
		}
	);

	const chartResults = $derived(dashboard?.response_chart ?? []);

	const monitors = $derived(dashboard?.monitors ?? []);

	const recentActivity = $derived(dashboard?.recent_activity ?? []);

	const chartMaximum = $derived.by(() => {
		if (chartResults.length === 0) {
			return 1;
		}

		return Math.max(...chartResults.map((item) => item.response_time_ms), 1);
	});

	const chartPoints = $derived.by(() => {
		if (chartResults.length === 0) {
			return '';
		}

		return chartResults
			.map((item, index) => {
				const x = getChartX(index, chartResults.length);

				const y = getChartY(item.response_time_ms, chartMaximum);

				return `${x},${y}`;
			})
			.join(' ');
	});

	const activeMonitorCount = $derived(
		summary.up_monitors + summary.down_monitors + summary.waiting_monitors
	);

	onMount(() => {
		void loadDashboard(true);
	});

	function getErrorMessage(error: unknown, fallback: string): string {
		if (error instanceof Error) {
			return error.message;
		}

		return fallback;
	}

	async function loadDashboard(initial = false): Promise<void> {
		if (initial) {
			loading = true;
		} else {
			refreshing = true;
		}

		try {
			dashboard = await getDashboard();
		} catch (error) {
			await showError(
				'Gagal memuat dashboard',
				getErrorMessage(error, 'Data dashboard tidak dapat dimuat.')
			);
		} finally {
			loading = false;
			refreshing = false;
		}
	}

	async function refreshDashboard(): Promise<void> {
		await loadDashboard(false);
	}

	async function openMonitors(): Promise<void> {
		await goto(resolve('/monitors'));
	}

	async function openHistory(): Promise<void> {
		await goto(resolve('/history'));
	}

	async function openMonitor(id: number): Promise<void> {
		await goto(resolve(`/monitors/${id}`));
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

	function getStatusClass(value: string): string {
		if (value === 'UP') {
			return 'bg-emerald-50 text-emerald-600 ring-emerald-100';
		}

		if (value === 'DOWN') {
			return 'bg-red-50 text-red-600 ring-red-100';
		}

		if (value === 'PAUSED') {
			return 'bg-amber-50 text-amber-600 ring-amber-100';
		}

		return 'bg-slate-100 text-slate-500 ring-slate-200';
	}

	function getActivityDotClass(value: string): string {
		if (value === 'UP') {
			return 'bg-emerald-500';
		}

		if (value === 'DOWN') {
			return 'bg-red-500';
		}

		return 'bg-slate-400';
	}

	function formatResponseTime(value: number): string {
		if (value <= 0) {
			return '-';
		}

		return `${value} ms`;
	}
</script>

<svelte:head>
	<title>Dashboard | GoWatch</title>
</svelte:head>

<PageHeader title="Dashboard" description="Pantau status website dan API dalam satu tempat." />

<div class="-mt-14 px-4 pb-10 sm:px-6 lg:px-8">
	<div class="mx-auto max-w-7xl">
		<div class="mb-5 flex justify-end">
			<div class="group relative">
				<button
					type="button"
					aria-label="Refresh dashboard"
					disabled={refreshing}
					onclick={() => {
						void refreshDashboard();
					}}
					class="flex h-10 w-10 items-center justify-center rounded-lg bg-white text-slate-500 shadow-lg transition hover:text-sky-500 disabled:cursor-not-allowed disabled:opacity-50"
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
				</button>

				<div
					class="pointer-events-none absolute right-0 bottom-full z-20 mb-2 rounded bg-slate-800 px-2 py-1 text-[10px] font-semibold whitespace-nowrap text-white opacity-0 shadow transition-opacity group-hover:opacity-100"
				>
					Refresh Dashboard
				</div>
			</div>
		</div>

		{#if loading}
			<div class="flex min-h-72 items-center justify-center rounded-lg bg-white shadow-lg">
				<div class="flex items-center gap-3">
					<div
						class="h-5 w-5 animate-spin rounded-full border-2 border-slate-200 border-t-sky-500"
					></div>

					<p class="text-sm text-slate-400">Memuat dashboard...</p>
				</div>
			</div>
		{:else}
			<div class="grid gap-5 sm:grid-cols-2 xl:grid-cols-4">
				<StatCard
					title="Total Monitor"
					value={summary.total_monitors}
					subtitle={`${activeMonitorCount} aktif • ${summary.paused_monitors} paused`}
					tone="sky"
					icon="monitor"
				/>

				<StatCard
					title="Service UP"
					value={summary.up_monitors}
					subtitle="Monitor aktif dengan status UP"
					tone="emerald"
					icon="up"
				/>

				<StatCard
					title="Service DOWN"
					value={summary.down_monitors}
					subtitle="Monitor aktif dengan status DOWN"
					tone="red"
					icon="down"
				/>

				<StatCard
					title="Uptime 24 Jam"
					value={`${summary.uptime_24h_percentage.toFixed(2)}%`}
					subtitle={`7 hari: ${summary.uptime_7d_percentage.toFixed(2)}%`}
					tone="orange"
					icon="performance"
				/>
			</div>

			<div class="mt-6 grid gap-6 xl:grid-cols-[minmax(0,1.6fr)_minmax(320px,0.7fr)]">
				<div class="overflow-hidden rounded-lg bg-slate-800 shadow-lg">
					<div
						class="flex flex-col gap-5 border-b border-slate-700 px-6 py-5 md:flex-row md:items-center md:justify-between"
					>
						<div>
							<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">
								Performance
							</p>

							<h2 class="mt-1 text-lg font-semibold text-white">Response Time — 24 Jam</h2>

							<p class="mt-1 text-xs text-slate-400">Hingga 48 pengecekan terbaru.</p>
						</div>

						<div class="grid grid-cols-2 gap-5">
							<div>
								<p class="text-[9px] font-bold tracking-wide text-slate-500 uppercase">Average</p>

								<p class="mt-1 text-sm font-semibold text-white">
									{formatResponseTime(summary.average_response_time_24h_ms)}
								</p>
							</div>

							<div>
								<p class="text-[9px] font-bold tracking-wide text-slate-500 uppercase">Checks</p>

								<p class="mt-1 text-sm font-semibold text-white">
									{summary.total_checks_24h}
								</p>
							</div>
						</div>
					</div>

					<div class="p-6">
						{#if chartResults.length === 0}
							<div
								class="flex h-72 items-center justify-center rounded-lg border border-dashed border-slate-600"
							>
								<div class="text-center">
									<p class="text-sm font-medium text-slate-300">Belum ada data 24 jam</p>

									<p class="mt-2 text-xs text-slate-500">
										Grafik akan muncul setelah worker melakukan pengecekan.
									</p>
								</div>
							</div>
						{:else}
							<div class="mb-4 flex justify-between text-[10px] font-medium text-slate-500">
								<span>
									{chartMaximum} ms
								</span>

								<span>
									{chartResults.length}
									check points
								</span>
							</div>

							<div class="h-72 w-full">
								<svg
									viewBox="0 0 1000 220"
									preserveAspectRatio="none"
									class="h-full w-full"
									role="img"
									aria-label="Grafik response time 24 jam"
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

									{#each chartResults as item, index (item.id)}
										<circle
											cx={getChartX(index, chartResults.length)}
											cy={getChartY(item.response_time_ms, chartMaximum)}
											r="4"
											fill="currentColor"
											vector-effect="non-scaling-stroke"
											class={item.status === 'DOWN' ? 'text-red-400' : 'text-sky-300'}
										>
											<title>
												{item.monitor_name} •
												{item.status} •
												{item.response_time_ms} ms
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

				<div class="rounded-lg bg-white p-6 shadow-lg">
					<div class="border-b border-slate-100 pb-5">
						<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">
							Monitor Status
						</p>

						<h2 class="mt-1 text-lg font-semibold text-slate-800">Service Distribution</h2>
					</div>

					<div class="mt-6 space-y-6">
						<div>
							<div class="mb-2 flex items-center justify-between">
								<div class="flex items-center gap-2">
									<span class="h-2.5 w-2.5 rounded-full bg-emerald-500"></span>

									<span class="text-xs font-semibold text-slate-600"> UP </span>
								</div>

								<span class="text-sm font-semibold text-emerald-600">
									{summary.up_monitors}
								</span>
							</div>

							<div class="h-2 overflow-hidden rounded-full bg-slate-100">
								<div
									class="h-full rounded-full bg-emerald-500"
									style={`width: ${
										summary.total_monitors > 0
											? (summary.up_monitors / summary.total_monitors) * 100
											: 0
									}%`}
								></div>
							</div>
						</div>

						<div>
							<div class="mb-2 flex items-center justify-between">
								<div class="flex items-center gap-2">
									<span class="h-2.5 w-2.5 rounded-full bg-red-500"></span>

									<span class="text-xs font-semibold text-slate-600"> DOWN </span>
								</div>

								<span class="text-sm font-semibold text-red-500">
									{summary.down_monitors}
								</span>
							</div>

							<div class="h-2 overflow-hidden rounded-full bg-slate-100">
								<div
									class="h-full rounded-full bg-red-500"
									style={`width: ${
										summary.total_monitors > 0
											? (summary.down_monitors / summary.total_monitors) * 100
											: 0
									}%`}
								></div>
							</div>
						</div>

						<div>
							<div class="mb-2 flex items-center justify-between">
								<div class="flex items-center gap-2">
									<span class="h-2.5 w-2.5 rounded-full bg-amber-500"></span>

									<span class="text-xs font-semibold text-slate-600"> PAUSED </span>
								</div>

								<span class="text-sm font-semibold text-amber-500">
									{summary.paused_monitors}
								</span>
							</div>

							<div class="h-2 overflow-hidden rounded-full bg-slate-100">
								<div
									class="h-full rounded-full bg-amber-500"
									style={`width: ${
										summary.total_monitors > 0
											? (summary.paused_monitors / summary.total_monitors) * 100
											: 0
									}%`}
								></div>
							</div>
						</div>

						<div>
							<div class="mb-2 flex items-center justify-between">
								<div class="flex items-center gap-2">
									<span class="h-2.5 w-2.5 rounded-full bg-slate-400"></span>

									<span class="text-xs font-semibold text-slate-600"> WAITING </span>
								</div>

								<span class="text-sm font-semibold text-slate-500">
									{summary.waiting_monitors}
								</span>
							</div>

							<div class="h-2 overflow-hidden rounded-full bg-slate-100">
								<div
									class="h-full rounded-full bg-slate-400"
									style={`width: ${
										summary.total_monitors > 0
											? (summary.waiting_monitors / summary.total_monitors) * 100
											: 0
									}%`}
								></div>
							</div>
						</div>
					</div>

					<div class="mt-7 grid grid-cols-2 gap-3 border-t border-slate-100 pt-6">
						<div class="rounded-lg bg-slate-50 p-4">
							<p class="text-[9px] font-bold tracking-wide text-slate-400 uppercase">Uptime 7D</p>

							<p class="mt-2 text-lg font-semibold text-slate-700">
								{summary.uptime_7d_percentage.toFixed(2)}%
							</p>
						</div>

						<div class="rounded-lg bg-slate-50 p-4">
							<p class="text-[9px] font-bold tracking-wide text-slate-400 uppercase">Checks 7D</p>

							<p class="mt-2 text-lg font-semibold text-slate-700">
								{summary.total_checks_7d}
							</p>
						</div>
					</div>
				</div>
			</div>

			<div class="mt-6 grid gap-6 xl:grid-cols-2">
				<div class="overflow-hidden rounded-lg bg-white shadow-lg">
					<div class="flex items-center justify-between border-b border-slate-100 px-6 py-5">
						<div>
							<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">Services</p>

							<h2 class="mt-1 text-base font-semibold text-slate-800">Monitor Status</h2>
						</div>

						<button
							type="button"
							onclick={() => {
								void openMonitors();
							}}
							class="text-xs font-semibold text-sky-500 transition hover:text-sky-600"
						>
							Lihat Semua
						</button>
					</div>

					{#if monitors.length === 0}
						<div class="flex min-h-52 items-center justify-center p-8 text-center">
							<div>
								<p class="text-sm font-semibold text-slate-700">Belum ada monitor</p>

								<p class="mt-2 text-xs text-slate-400">
									Tambahkan website atau API untuk mulai monitoring.
								</p>
							</div>
						</div>
					{:else}
						<div class="divide-y divide-slate-100">
							{#each monitors as monitor (monitor.id)}
								<button
									type="button"
									onclick={() => {
										void openMonitor(monitor.id);
									}}
									class="flex w-full items-center justify-between gap-4 px-6 py-4 text-left transition hover:bg-slate-50"
								>
									<div class="min-w-0">
										<div class="flex items-center gap-2">
											<p class="truncate text-sm font-semibold text-slate-700">
												{monitor.name}
											</p>

											<span
												class={`inline-flex rounded-full px-2 py-0.5 text-[9px] font-bold ring-1 ring-inset ${getStatusClass(
													monitor.current_status
												)}`}
											>
												{monitor.current_status}
											</span>
										</div>

										<p class="mt-1 max-w-sm truncate text-xs text-slate-400">
											{monitor.url}
										</p>
									</div>

									<div class="shrink-0 text-right">
										<p class="text-xs font-semibold text-slate-600">
											{formatResponseTime(monitor.last_response_time_ms)}
										</p>

										<p class="mt-1 text-[10px] text-slate-400">
											{monitor.last_checked_at
												? formatDate(monitor.last_checked_at)
												: 'Belum dicek'}
										</p>
									</div>
								</button>
							{/each}
						</div>
					{/if}
				</div>

				<div class="overflow-hidden rounded-lg bg-white shadow-lg">
					<div class="flex items-center justify-between border-b border-slate-100 px-6 py-5">
						<div>
							<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">Activity</p>

							<h2 class="mt-1 text-base font-semibold text-slate-800">Recent Checks</h2>
						</div>

						<button
							type="button"
							onclick={() => {
								void openHistory();
							}}
							class="text-xs font-semibold text-sky-500 transition hover:text-sky-600"
						>
							Lihat History
						</button>
					</div>

					{#if recentActivity.length === 0}
						<div class="flex min-h-52 items-center justify-center p-8 text-center">
							<p class="text-sm text-slate-400">Belum ada aktivitas monitoring.</p>
						</div>
					{:else}
						<div class="divide-y divide-slate-100">
							{#each recentActivity as activity (activity.id)}
								<button
									type="button"
									onclick={() => {
										void openMonitor(activity.monitor_id);
									}}
									class="flex w-full items-center gap-4 px-6 py-4 text-left transition hover:bg-slate-50"
								>
									<span
										class={`h-2.5 w-2.5 shrink-0 rounded-full ${getActivityDotClass(
											activity.status
										)}`}
									></span>

									<div class="min-w-0 flex-1">
										<div class="flex items-center justify-between gap-4">
											<p class="truncate text-xs font-semibold text-slate-700">
												{activity.monitor_name}
											</p>

											<span class="shrink-0 text-[10px] text-slate-400">
												{formatDate(activity.checked_at)}
											</span>
										</div>

										<div class="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-[10px]">
											<span
												class={activity.status === 'UP'
													? 'font-semibold text-emerald-600'
													: 'font-semibold text-red-500'}
											>
												{activity.status}
											</span>

											<span class="text-slate-400">
												HTTP
												{activity.http_status || '-'}
											</span>

											<span class="text-slate-400">
												{formatResponseTime(activity.response_time_ms)}
											</span>
										</div>

										{#if activity.error}
											<p class="mt-1 truncate text-[10px] text-red-400">
												{activity.error}
											</p>
										{/if}
									</div>
								</button>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		{/if}
	</div>
</div>
