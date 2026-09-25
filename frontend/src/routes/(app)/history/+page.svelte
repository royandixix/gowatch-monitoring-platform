<script lang="ts">
	import { onMount } from 'svelte';

	import { goto } from '$app/navigation';

	import { resolve } from '$app/paths';

	import { getHistory } from '$lib/api/history';

	import { getMonitors } from '$lib/api/monitors';

	import PageHeader from '$lib/components/layout/PageHeader.svelte';

	import type { HistoryStatus, MonitorResultRecord } from '$lib/types/history';

	import type { Monitor } from '$lib/types/monitor';

	import { showError } from '$lib/utils/alert';

	import { formatDate } from '$lib/utils/format';

	let results = $state<MonitorResultRecord[]>([]);

	let monitors = $state<Monitor[]>([]);

	let loading = $state(true);

	let status = $state<HistoryStatus>('ALL');

	let selectedMonitor = $state('');

	let pageNumber = $state(1);

	let perPage = $state(20);

	let total = $state(0);

	let totalPages = $state(0);

	const showingFrom = $derived(total === 0 ? 0 : (pageNumber - 1) * perPage + 1);

	const showingTo = $derived(Math.min(pageNumber * perPage, total));

	const displayedTotalPages = $derived(Math.max(totalPages, 1));

	const upCount = $derived(results.filter((result) => result.status === 'UP').length);

	const downCount = $derived(results.filter((result) => result.status === 'DOWN').length);

	onMount(() => {
		void initialize();
	});

	async function initialize(): Promise<void> {
		loading = true;

		try {
			const monitorResponse = await getMonitors();

			monitors = monitorResponse.monitors;

			await loadHistory();
		} catch (error) {
			await showError(
				'Gagal memuat history',
				getErrorMessage(error, 'Data history monitoring tidak dapat dimuat.')
			);

			loading = false;
		}
	}

	function getErrorMessage(error: unknown, fallback: string): string {
		if (error instanceof Error) {
			return error.message;
		}

		return fallback;
	}

	async function loadHistory(): Promise<void> {
		loading = true;

		try {
			const monitorID = selectedMonitor ? Number(selectedMonitor) : undefined;

			const response = await getHistory({
				status,
				monitorId: monitorID,
				page: pageNumber,
				perPage
			});

			results = response.results;

			total = response.total;

			totalPages = response.total_pages;

			if (totalPages > 0 && pageNumber > totalPages) {
				pageNumber = totalPages;

				await loadHistory();

				return;
			}
		} catch (error) {
			await showError(
				'Gagal memuat history',
				getErrorMessage(error, 'Terjadi kesalahan saat mengambil history monitoring.')
			);
		} finally {
			loading = false;
		}
	}

	async function setStatus(nextStatus: HistoryStatus): Promise<void> {
		if (status === nextStatus) {
			return;
		}

		status = nextStatus;

		pageNumber = 1;

		await loadHistory();
	}

	async function changeMonitor(): Promise<void> {
		pageNumber = 1;

		await loadHistory();
	}

	async function changePerPage(): Promise<void> {
		pageNumber = 1;

		await loadHistory();
	}

	async function refreshHistory(): Promise<void> {
		await loadHistory();
	}

	async function goToPage(nextPage: number): Promise<void> {
		if (nextPage < 1 || nextPage > displayedTotalPages || nextPage === pageNumber) {
			return;
		}

		pageNumber = nextPage;

		await loadHistory();

		window.scrollTo({
			top: 0,
			behavior: 'smooth'
		});
	}

	async function openMonitor(monitorID: number): Promise<void> {
		await goto(resolve(`/monitors/${monitorID}`));
	}

	function getStatusClass(value: string): string {
		if (value === 'UP') {
			return 'bg-emerald-50 text-emerald-600 ring-emerald-100';
		}

		if (value === 'DOWN') {
			return 'bg-red-50 text-red-600 ring-red-100';
		}

		return 'bg-slate-100 text-slate-500 ring-slate-200';
	}

	function getHTTPStatusClass(value: number): string {
		if (value >= 200 && value < 300) {
			return 'text-emerald-600';
		}

		if (value >= 400) {
			return 'text-red-500';
		}

		return 'text-slate-500';
	}

	function formatResponseTime(value: number): string {
		if (value <= 0) {
			return '-';
		}

		return `${value} ms`;
	}
</script>

<svelte:head>
	<title>History | GoWatch</title>
</svelte:head>

<PageHeader
	title="History"
	description="Lihat riwayat hasil pengecekan website dan API yang dipantau GoWatch."
/>

<div class="-mt-14 px-4 pb-10 sm:px-6 lg:px-8">
	<div class="mx-auto max-w-7xl">
		<div class="grid gap-4 sm:grid-cols-3">
			<div class="rounded-lg bg-white p-5 shadow-lg">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">
							Total History
						</p>

						<p class="mt-2 text-2xl font-semibold text-slate-700">
							{total}
						</p>
					</div>

					<div
						class="flex h-11 w-11 items-center justify-center rounded-full bg-sky-50 text-sky-500"
					>
						<svg
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="1.8"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="h-5 w-5"
							aria-hidden="true"
						>
							<path d="M3 12a9 9 0 1 0 3-6.7"></path>

							<path d="M3 4v5h5"></path>

							<path d="M12 7v5l3 2"></path>
						</svg>
					</div>
				</div>
			</div>

			<div class="rounded-lg bg-white p-5 shadow-lg">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">
							UP di Halaman
						</p>

						<p class="mt-2 text-2xl font-semibold text-emerald-600">
							{upCount}
						</p>
					</div>

					<div
						class="flex h-11 w-11 items-center justify-center rounded-full bg-emerald-50 text-emerald-500"
					>
						<svg
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="h-5 w-5"
							aria-hidden="true"
						>
							<path d="m5 12 4 4L19 6"></path>
						</svg>
					</div>
				</div>
			</div>

			<div class="rounded-lg bg-white p-5 shadow-lg">
				<div class="flex items-center justify-between">
					<div>
						<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">
							DOWN di Halaman
						</p>

						<p class="mt-2 text-2xl font-semibold text-red-500">
							{downCount}
						</p>
					</div>

					<div
						class="flex h-11 w-11 items-center justify-center rounded-full bg-red-50 text-red-500"
					>
						<svg
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="h-5 w-5"
							aria-hidden="true"
						>
							<path d="M6 6l12 12"></path>

							<path d="M18 6 6 18"></path>
						</svg>
					</div>
				</div>
			</div>
		</div>

		<div class="mt-6 rounded-lg bg-white p-5 shadow-lg">
			<div class="flex flex-col gap-5 xl:flex-row xl:items-end xl:justify-between">
				<div>
					<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">Status</p>

					<div class="mt-2 flex flex-wrap gap-2">
						<button
							type="button"
							aria-pressed={status === 'ALL'}
							onclick={() => {
								void setStatus('ALL');
							}}
							class={`rounded-lg px-4 py-2 text-xs font-semibold transition ${
								status === 'ALL'
									? 'bg-slate-800 text-white shadow'
									: 'bg-slate-100 text-slate-500 hover:bg-slate-200'
							}`}
						>
							Semua
						</button>

						<button
							type="button"
							aria-pressed={status === 'UP'}
							onclick={() => {
								void setStatus('UP');
							}}
							class={`rounded-lg px-4 py-2 text-xs font-semibold transition ${
								status === 'UP'
									? 'bg-emerald-500 text-white shadow'
									: 'bg-emerald-50 text-emerald-600 hover:bg-emerald-100'
							}`}
						>
							UP
						</button>

						<button
							type="button"
							aria-pressed={status === 'DOWN'}
							onclick={() => {
								void setStatus('DOWN');
							}}
							class={`rounded-lg px-4 py-2 text-xs font-semibold transition ${
								status === 'DOWN'
									? 'bg-red-500 text-white shadow'
									: 'bg-red-50 text-red-600 hover:bg-red-100'
							}`}
						>
							DOWN
						</button>
					</div>
				</div>

				<div class="grid gap-4 sm:grid-cols-2 xl:flex xl:items-end">
					<div class="min-w-[220px]">
						<label
							for="history-monitor"
							class="mb-2 block text-[10px] font-bold tracking-wide text-slate-400 uppercase"
						>
							Monitor
						</label>

						<select
							id="history-monitor"
							bind:value={selectedMonitor}
							onchange={() => {
								void changeMonitor();
							}}
							class="h-10 w-full rounded-lg border border-slate-200 bg-white px-3 text-xs font-medium text-slate-600 transition outline-none focus:border-sky-400 focus:ring-4 focus:ring-sky-50"
						>
							<option value=""> Semua monitor </option>

							{#each monitors as monitor (monitor.id)}
								<option value={String(monitor.id)}>
									{monitor.name}
								</option>
							{/each}
						</select>
					</div>

					<div>
						<label
							for="history-limit"
							class="mb-2 block text-[10px] font-bold tracking-wide text-slate-400 uppercase"
						>
							Data / Halaman
						</label>

						<select
							id="history-limit"
							bind:value={perPage}
							onchange={() => {
								void changePerPage();
							}}
							class="h-10 w-full min-w-[130px] rounded-lg border border-slate-200 bg-white px-3 text-xs font-medium text-slate-600 transition outline-none focus:border-sky-400 focus:ring-4 focus:ring-sky-50"
						>
							<option value={20}> 20 </option>

							<option value={50}> 50 </option>

							<option value={100}> 100 </option>
						</select>
					</div>

					<div class="group relative">
						<button
							type="button"
							aria-label="Refresh history"
							onclick={() => {
								void refreshHistory();
							}}
							disabled={loading}
							class="flex h-10 w-10 items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-500 transition hover:bg-slate-50 hover:text-sky-500 disabled:cursor-not-allowed disabled:opacity-50"
						>
							<svg
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="1.8"
								stroke-linecap="round"
								stroke-linejoin="round"
								class={`h-4 w-4 ${loading ? 'animate-spin' : ''}`}
								aria-hidden="true"
							>
								<path d="M20 11a8.1 8.1 0 0 0-15.5-2M4 4v5h5"></path>

								<path d="M4 13a8.1 8.1 0 0 0 15.5 2M20 20v-5h-5"></path>
							</svg>
						</button>

						<div
							class="pointer-events-none absolute right-0 bottom-full z-20 mb-2 rounded bg-slate-800 px-2 py-1 text-[10px] font-medium whitespace-nowrap text-white opacity-0 shadow transition-opacity group-hover:opacity-100"
						>
							Refresh
						</div>
					</div>
				</div>
			</div>
		</div>

		<div class="mt-6 overflow-hidden rounded-lg bg-white shadow-lg">
			<div
				class="flex flex-col gap-2 border-b border-slate-100 px-5 py-4 sm:flex-row sm:items-center sm:justify-between"
			>
				<div>
					<h2 class="text-sm font-semibold text-slate-800">Riwayat Pengecekan</h2>

					<p class="mt-1 text-xs text-slate-400">
						Menampilkan {showingFrom}–{showingTo}
						dari {total} hasil.
					</p>
				</div>

				<p class="text-xs font-medium text-slate-400">
					Halaman {pageNumber}
					dari {displayedTotalPages}
				</p>
			</div>

			{#if loading}
				<div class="flex min-h-64 items-center justify-center">
					<div class="flex items-center gap-3">
						<div
							class="h-5 w-5 animate-spin rounded-full border-2 border-slate-200 border-t-sky-500"
						></div>

						<p class="text-sm text-slate-400">Memuat history...</p>
					</div>
				</div>
			{:else if results.length === 0}
				<div class="flex min-h-64 flex-col items-center justify-center p-8 text-center">
					<div
						class="flex h-12 w-12 items-center justify-center rounded-full bg-slate-100 text-slate-400"
					>
						<svg
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="1.8"
							stroke-linecap="round"
							stroke-linejoin="round"
							class="h-6 w-6"
							aria-hidden="true"
						>
							<path d="M3 12a9 9 0 1 0 3-6.7"></path>

							<path d="M3 4v5h5"></path>
						</svg>
					</div>

					<p class="mt-4 text-sm font-semibold text-slate-700">History tidak ditemukan</p>

					<p class="mt-2 max-w-md text-xs leading-5 text-slate-400">
						Belum ada hasil pengecekan yang sesuai dengan filter yang dipilih.
					</p>
				</div>
			{:else}
				<div class="overflow-x-auto">
					<table class="w-full min-w-[980px]">
						<thead class="bg-slate-50">
							<tr>
								<th
									class="px-5 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
								>
									Monitor
								</th>

								<th
									class="px-5 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
								>
									Status
								</th>

								<th
									class="px-5 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
								>
									HTTP
								</th>

								<th
									class="px-5 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
								>
									Response
								</th>

								<th
									class="px-5 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
								>
									Waktu
								</th>

								<th
									class="px-5 py-4 text-left text-[10px] font-bold tracking-wide text-slate-400 uppercase"
								>
									Error
								</th>

								<th
									class="px-5 py-4 text-right text-[10px] font-bold tracking-wide text-slate-400 uppercase"
								>
									Detail
								</th>
							</tr>
						</thead>

						<tbody class="divide-y divide-slate-100">
							{#each results as result (result.id)}
								<tr class="transition hover:bg-slate-50">
									<td class="px-5 py-4">
										<p class="max-w-[240px] truncate text-sm font-semibold text-slate-700">
											{result.monitor_name}
										</p>

										<p class="mt-1 max-w-[240px] truncate text-xs text-slate-400">
											{result.monitor_url}
										</p>
									</td>

									<td class="px-5 py-4">
										<span
											class={`inline-flex rounded-full px-2.5 py-1 text-[10px] font-bold ring-1 ring-inset ${getStatusClass(
												result.status
											)}`}
										>
											{result.status}
										</span>
									</td>

									<td class="px-5 py-4">
										<span class={`text-xs font-semibold ${getHTTPStatusClass(result.http_status)}`}>
											{result.http_status > 0 ? result.http_status : '-'}
										</span>
									</td>

									<td class="px-5 py-4 text-xs font-medium text-slate-600">
										{formatResponseTime(result.response_time_ms)}
									</td>

									<td class="px-5 py-4 text-xs text-slate-500">
										{formatDate(result.checked_at)}
									</td>

									<td class="px-5 py-4">
										<p
											class={`max-w-[220px] truncate text-xs ${
												result.error ? 'text-red-500' : 'text-slate-300'
											}`}
											title={result.error || undefined}
										>
											{result.error || '-'}
										</p>
									</td>

									<td class="px-5 py-4">
										<div class="flex justify-end">
											<div class="group relative">
												<button
													type="button"
													aria-label={`Lihat ${result.monitor_name}`}
													onclick={() => {
														void openMonitor(result.monitor_id);
													}}
													class="flex h-8 w-8 items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-500 transition hover:bg-slate-100 hover:text-slate-800"
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
														<path
															d="M2.25 12s3.75-6.75 9.75-6.75S21.75 12 21.75 12 18 18.75 12 18.75 2.25 12 2.25 12Z"
														></path>

														<circle cx="12" cy="12" r="3"></circle>
													</svg>
												</button>

												<div
													class="pointer-events-none absolute right-0 bottom-full z-20 mb-2 rounded bg-slate-800 px-2 py-1 text-[10px] font-medium whitespace-nowrap text-white opacity-0 shadow transition-opacity group-hover:opacity-100"
												>
													Detail Monitor
												</div>
											</div>
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}

			{#if !loading && total > 0}
				<div
					class="flex flex-col gap-4 border-t border-slate-100 px-5 py-4 sm:flex-row sm:items-center sm:justify-between"
				>
					<p class="text-xs text-slate-400">
						{showingFrom}–{showingTo}
						dari {total} data
					</p>

					<div class="flex items-center gap-2">
						<button
							type="button"
							disabled={pageNumber <= 1}
							onclick={() => {
								void goToPage(pageNumber - 1);
							}}
							class="flex h-9 items-center gap-2 rounded-lg border border-slate-200 bg-white px-3 text-xs font-semibold text-slate-500 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
						>
							<svg
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="h-3.5 w-3.5"
								aria-hidden="true"
							>
								<path d="m15 18-6-6 6-6"></path>
							</svg>

							Sebelumnya
						</button>

						<div
							class="flex h-9 min-w-24 items-center justify-center rounded-lg bg-slate-100 px-3 text-xs font-semibold text-slate-600"
						>
							{pageNumber}
							/
							{displayedTotalPages}
						</div>

						<button
							type="button"
							disabled={pageNumber >= displayedTotalPages}
							onclick={() => {
								void goToPage(pageNumber + 1);
							}}
							class="flex h-9 items-center gap-2 rounded-lg border border-slate-200 bg-white px-3 text-xs font-semibold text-slate-500 transition hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
						>
							Selanjutnya

							<svg
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
								class="h-3.5 w-3.5"
								aria-hidden="true"
							>
								<path d="m9 18 6-6-6-6"></path>
							</svg>
						</button>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
