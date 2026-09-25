<script lang="ts">
	import { onMount } from 'svelte';

	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';

	import type { CreateMonitorPayload, Monitor } from '$lib/types/monitor';

	let {
		monitor = null,
		loading = false,
		onSubmit,
		onCancel
	}: {
		monitor?: Monitor | null;
		loading?: boolean;
		onSubmit: (payload: CreateMonitorPayload) => Promise<void> | void;
		onCancel?: () => void;
	} = $props();

	let name = $state('');
	let url = $state('');
	let monitorType = $state<'website' | 'api'>('website');
	let intervalSeconds = $state('60');
	let active = $state(true);

	let nameError = $state('');
	let urlError = $state('');
	let intervalError = $state('');

	onMount(() => {
		if (!monitor) {
			return;
		}

		name = monitor.name;
		url = monitor.url;
		monitorType = monitor.monitor_type;
		intervalSeconds = String(monitor.interval_seconds);
		active = monitor.is_active;
	});

	function validate(): boolean {
		nameError = '';
		urlError = '';
		intervalError = '';

		if (!name.trim()) {
			nameError = 'Nama monitor wajib diisi.';
		}

		if (!url.trim()) {
			urlError = 'URL wajib diisi.';
		} else {
			try {
				const parsed = new URL(url.trim());

				if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') {
					urlError = 'Gunakan URL http atau https.';
				}
			} catch {
				urlError = 'Format URL tidak valid.';
			}
		}

		const interval = Number(intervalSeconds);

		if (!Number.isFinite(interval) || interval < 30) {
			intervalError = 'Interval minimal 30 detik.';
		}

		return !nameError && !urlError && !intervalError;
	}

	async function submit(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		if (!validate()) {
			return;
		}

		await onSubmit({
			name: name.trim(),
			url: url.trim(),
			monitor_type: monitorType,
			interval_seconds: Number(intervalSeconds),
			is_active: active
		});
	}
</script>

<form onsubmit={submit} class="rounded-md bg-white p-6 shadow-lg">
	<div class="mb-6">
		<h2 class="text-base font-semibold text-slate-800">
			{monitor ? 'Edit Monitor' : 'Tambah Monitor'}
		</h2>

		<p class="mt-1 text-xs leading-5 text-slate-400">
			{monitor
				? 'Perbarui konfigurasi monitor.'
				: 'Tambahkan website atau API yang ingin dipantau.'}
		</p>
	</div>

	<div class="space-y-5">
		<Input
			id="monitor-name"
			label="Nama monitor"
			placeholder="Contoh: Website Production"
			required
			bind:value={name}
			error={nameError}
		/>

		<Input
			id="monitor-url"
			label="URL"
			type="url"
			placeholder="https://example.com"
			required
			bind:value={url}
			error={urlError}
		/>

		<div class="grid gap-5 sm:grid-cols-2">
			<div>
				<label for="monitor-type" class="mb-2 block text-xs font-semibold text-slate-600">
					Tipe monitor
				</label>

				<select
					id="monitor-type"
					bind:value={monitorType}
					class="h-11 w-full rounded-md border border-slate-200 bg-white px-3 text-sm text-slate-700 outline-none focus:border-sky-400 focus:ring-4 focus:ring-sky-50"
				>
					<option value="website"> Website </option>

					<option value="api"> API </option>
				</select>
			</div>

			<div>
				<label for="monitor-interval" class="mb-2 block text-xs font-semibold text-slate-600">
					Interval
				</label>

				<select
					id="monitor-interval"
					bind:value={intervalSeconds}
					class={`h-11 w-full rounded-md border bg-white px-3 text-sm text-slate-700 outline-none ${
						intervalError
							? 'border-red-300'
							: 'border-slate-200 focus:border-sky-400 focus:ring-4 focus:ring-sky-50'
					}`}
				>
					<option value="30"> 30 detik </option>

					<option value="60"> 1 menit </option>

					<option value="120"> 2 menit </option>

					<option value="300"> 5 menit </option>

					<option value="600"> 10 menit </option>
				</select>

				{#if intervalError}
					<p class="mt-1.5 text-xs text-red-500">
						{intervalError}
					</p>
				{/if}
			</div>
		</div>

		<label
			class="flex cursor-pointer items-center justify-between rounded-md border border-slate-200 p-4"
		>
			<div>
				<p class="text-sm font-semibold text-slate-700">Monitoring aktif</p>

				<p class="mt-1 text-xs text-slate-400">
					Background worker akan memeriksa layanan secara otomatis.
				</p>
			</div>

			<input type="checkbox" bind:checked={active} class="h-4 w-4 accent-sky-500" />
		</label>
	</div>

	<div class="mt-6 flex justify-end gap-3">
		{#if onCancel}
			<Button variant="secondary" onclick={onCancel}>Batal</Button>
		{/if}

		<Button type="submit" {loading}>
			{monitor ? 'Simpan Perubahan' : 'Tambah Monitor'}
		</Button>
	</div>
</form>
