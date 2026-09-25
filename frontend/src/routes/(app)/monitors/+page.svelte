<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';

	import Swal from 'sweetalert2';
	import 'sweetalert2/dist/sweetalert2.min.css';

	import PageHeader from '$lib/components/layout/PageHeader.svelte';

	import MonitorCard from '$lib/components/monitor/MonitorCard.svelte';

	import MonitorForm from '$lib/components/monitor/MonitorForm.svelte';

	import MonitorTable from '$lib/components/monitor/MonitorTable.svelte';

	import { createMonitor, deleteMonitor, getMonitors, updateMonitor } from '$lib/api/monitors';

	import type { CreateMonitorPayload, Monitor } from '$lib/types/monitor';

	let monitors = $state<Monitor[]>([]);

	let loading = $state(true);

	let saving = $state(false);

	let showForm = $state(false);

	let editingMonitor = $state<Monitor | null>(null);

	const query = $derived(page.url.searchParams.get('search')?.trim().toLowerCase() ?? '');

	const filteredMonitors = $derived.by(() => {
		if (!query) {
			return monitors;
		}

		return monitors.filter(
			(monitor) =>
				monitor.name.toLowerCase().includes(query) ||
				monitor.url.toLowerCase().includes(query) ||
				monitor.monitor_type.toLowerCase().includes(query)
		);
	});

	onMount(() => {
		void loadMonitors();
	});

	function getErrorMessage(error: unknown, fallback: string): string {
		if (error instanceof Error) {
			return error.message;
		}

		return fallback;
	}

	async function showSuccess(title: string, message: string): Promise<void> {
		await Swal.fire({
			toast: true,
			position: 'top-end',
			icon: 'success',
			title,
			text: message,
			showConfirmButton: false,
			timer: 2600,
			timerProgressBar: true,
			showClass: {
				popup: 'swal2-show'
			},
			hideClass: {
				popup: 'swal2-hide'
			},
			customClass: {
				popup: '!rounded-xl !border !border-slate-200 !bg-white !shadow-xl',
				title: '!text-sm !font-semibold !text-slate-800',
				htmlContainer: '!text-xs !text-slate-500',
				timerProgressBar: '!bg-sky-500'
			}
		});
	}

	async function showError(title: string, message: string): Promise<void> {
		await Swal.fire({
			icon: 'error',
			title,
			text: message,
			confirmButtonText: 'Tutup',
			buttonsStyling: false,
			customClass: {
				popup: '!w-[min(92vw,420px)] !rounded-xl !p-7 !shadow-2xl',
				title: '!text-xl !font-semibold !text-slate-800',
				htmlContainer: '!mt-2 !text-sm !leading-6 !text-slate-500',
				actions: '!mt-6',
				confirmButton:
					'!h-10 !rounded-lg !bg-sky-500 !px-6 !text-xs !font-semibold !text-white !transition hover:!bg-sky-600'
			}
		});
	}

	async function confirmDelete(monitor: Monitor): Promise<boolean> {
		const result = await Swal.fire({
			icon: 'warning',
			title: 'Hapus monitor?',
			html: `
				<div style="font-size:13px;line-height:1.7;color:#64748b;">
					Monitor
					<strong style="color:#334155;">
						${escapeHtml(monitor.name)}
					</strong>
					akan dihapus beserta riwayat monitoring yang terkait.
				</div>
			`,
			showCancelButton: true,
			confirmButtonText: 'Hapus',
			cancelButtonText: 'Batal',
			reverseButtons: true,
			focusCancel: true,
			buttonsStyling: false,
			customClass: {
				popup: '!w-[min(92vw,430px)] !rounded-xl !p-7 !shadow-2xl',
				title: '!text-xl !font-semibold !text-slate-800',
				htmlContainer: '!mt-2',
				actions: '!mt-6 !gap-2',
				confirmButton:
					'!h-10 !rounded-lg !bg-red-500 !px-6 !text-xs !font-semibold !text-white !transition hover:!bg-red-600',
				cancelButton:
					'!h-10 !rounded-lg !border !border-slate-200 !bg-white !px-6 !text-xs !font-semibold !text-slate-600 !transition hover:!bg-slate-50'
			}
		});

		return result.isConfirmed;
	}

	function escapeHtml(value: string): string {
		return value
			.replaceAll('&', '&amp;')
			.replaceAll('<', '&lt;')
			.replaceAll('>', '&gt;')
			.replaceAll('"', '&quot;')
			.replaceAll("'", '&#039;');
	}

	async function loadMonitors(): Promise<void> {
		loading = true;

		try {
			const response = await getMonitors();

			monitors = response.monitors;
		} catch (error) {
			await showError(
				'Gagal memuat monitor',
				getErrorMessage(error, 'Data monitor tidak dapat dimuat.')
			);
		} finally {
			loading = false;
		}
	}

	function openCreate(): void {
		editingMonitor = null;
		showForm = true;
	}

	function openEdit(monitor: Monitor): void {
		editingMonitor = monitor;

		showForm = true;

		window.scrollTo({
			top: 0,
			behavior: 'smooth'
		});
	}

	function closeForm(): void {
		showForm = false;

		editingMonitor = null;
	}

	async function saveMonitor(payload: CreateMonitorPayload): Promise<void> {
		if (saving) {
			return;
		}

		saving = true;

		try {
			if (editingMonitor) {
				const monitorName = editingMonitor.name;

				await updateMonitor(editingMonitor.id, payload);

				closeForm();

				await loadMonitors();

				await showSuccess('Monitor diperbarui', `Monitor "${monitorName}" berhasil diperbarui.`);

				return;
			}

			await createMonitor(payload);

			closeForm();

			await loadMonitors();

			await showSuccess(
				'Monitor ditambahkan',
				'Monitor baru berhasil ditambahkan dan siap dipantau.'
			);
		} catch (error) {
			await showError(
				'Gagal menyimpan monitor',
				getErrorMessage(error, 'Terjadi kesalahan saat menyimpan monitor.')
			);
		} finally {
			saving = false;
		}
	}

	async function removeMonitor(monitor: Monitor): Promise<void> {
		const confirmed = await confirmDelete(monitor);

		if (!confirmed) {
			return;
		}

		try {
			Swal.fire({
				title: 'Menghapus monitor...',
				text: 'Mohon tunggu sebentar.',
				allowOutsideClick: false,
				allowEscapeKey: false,
				showConfirmButton: false,
				customClass: {
					popup: '!w-[min(92vw,400px)] !rounded-xl !p-7 !shadow-2xl',
					title: '!text-lg !font-semibold !text-slate-800',
					htmlContainer: '!text-sm !text-slate-500'
				},
				didOpen: () => {
					Swal.showLoading();
				}
			});

			await deleteMonitor(monitor.id);

			await loadMonitors();

			Swal.close();

			await showSuccess('Monitor dihapus', `Monitor "${monitor.name}" berhasil dihapus.`);
		} catch (error) {
			Swal.close();

			await showError(
				'Gagal menghapus monitor',
				getErrorMessage(error, 'Terjadi kesalahan saat menghapus monitor.')
			);
		}
	}
</script>

<svelte:head>
	<title>Monitors | GoWatch</title>
</svelte:head>

<PageHeader title="Monitors" description="Kelola website dan API yang dipantau oleh GoWatch." />

<div class="-mt-14 px-4 pb-10 sm:px-6 lg:px-8">
	<div class="mx-auto max-w-7xl">
		<div class="mb-5 flex items-end justify-between gap-4">
			<div class="rounded-lg bg-white px-5 py-4 shadow-lg">
				<p class="text-[10px] font-bold tracking-wide text-slate-400 uppercase">Total Monitor</p>

				<p class="mt-1 text-2xl font-semibold text-slate-700">
					{monitors.length}
				</p>
			</div>

			<div class="group relative">
				<button
					type="button"
					aria-label="Tambah monitor"
					onclick={openCreate}
					class="flex h-11 w-11 items-center justify-center rounded-lg bg-sky-500 text-white shadow-lg shadow-sky-500/20 transition duration-200 hover:bg-sky-600 hover:shadow-sky-500/30"
				>
					<svg
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						class="h-5 w-5 transition-transform duration-200 group-hover:rotate-90"
						aria-hidden="true"
					>
						<path d="M12 5v14"></path>

						<path d="M5 12h14"></path>
					</svg>
				</button>

				<div
					class="pointer-events-none absolute right-0 bottom-full z-20 mb-2 rounded bg-slate-800 px-2.5 py-1.5 text-[10px] font-medium whitespace-nowrap text-white opacity-0 shadow-lg transition-opacity group-hover:opacity-100"
				>
					Tambah Monitor
				</div>
			</div>
		</div>

		{#if showForm}
			<div class="mb-6">
				{#key editingMonitor?.id ?? 'new'}
					<MonitorForm
						monitor={editingMonitor}
						loading={saving}
						onSubmit={saveMonitor}
						onCancel={closeForm}
					/>
				{/key}
			</div>
		{/if}

		{#if query}
			<div
				class="mb-4 flex items-center gap-2 rounded-lg border border-slate-100 bg-white px-4 py-3 text-xs text-slate-500 shadow-sm"
			>
				<svg
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					class="h-4 w-4 text-slate-400"
					aria-hidden="true"
				>
					<circle cx="11" cy="11" r="8"></circle>

					<path d="m21 21-4.3-4.3"></path>
				</svg>

				<span> Hasil pencarian: </span>

				<strong class="text-slate-700">
					{page.url.searchParams.get('search')}
				</strong>
			</div>
		{/if}

		{#if loading}
			<div class="flex min-h-48 items-center justify-center rounded-lg bg-white shadow-lg">
				<div class="flex items-center gap-3">
					<div
						class="h-5 w-5 animate-spin rounded-full border-2 border-slate-200 border-t-sky-500"
					></div>

					<p class="text-sm text-slate-400">Memuat monitor...</p>
				</div>
			</div>
		{:else if filteredMonitors.length === 0}
			<div class="rounded-lg bg-white p-12 text-center shadow-lg">
				<div
					class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-slate-100 text-slate-400"
				>
					<svg
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.8"
						class="h-6 w-6"
						aria-hidden="true"
					>
						<rect x="3" y="4" width="18" height="13" rx="2"></rect>

						<path d="M8 21h8"></path>

						<path d="M12 17v4"></path>
					</svg>
				</div>

				<p class="mt-4 text-base font-semibold text-slate-700">Monitor tidak ditemukan</p>

				<p class="mt-2 text-sm text-slate-400">
					{#if query}
						Tidak ada monitor yang sesuai dengan pencarian.
					{:else}
						Tambahkan website atau API yang ingin dipantau.
					{/if}
				</p>

				{#if !query}
					<button
						type="button"
						onclick={openCreate}
						class="mx-auto mt-5 flex h-10 items-center justify-center gap-2 rounded-lg bg-sky-500 px-4 text-xs font-semibold text-white transition hover:bg-sky-600"
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
							<path d="M12 5v14"></path>

							<path d="M5 12h14"></path>
						</svg>

						Tambah Monitor
					</button>
				{/if}
			</div>
		{:else}
			<div class="hidden md:block">
				<MonitorTable monitors={filteredMonitors} onEdit={openEdit} onDelete={removeMonitor} />
			</div>

			<div class="grid gap-4 md:hidden">
				{#each filteredMonitors as monitor (monitor.id)}
					<MonitorCard {monitor} onEdit={openEdit} onDelete={removeMonitor} />
				{/each}
			</div>
		{/if}
	</div>
</div>
