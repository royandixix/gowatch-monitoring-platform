<script lang="ts">
	import { fade, scale } from 'svelte/transition';

	type ModalType = 'info' | 'success' | 'error';

	interface Props {
		open?: boolean;
		title?: string;
		message?: string;
		type?: ModalType;
		buttonText?: string;
		onClose: () => void;
	}

	let {
		open = false,
		title = '',
		message = '',
		type = 'info',
		buttonText = 'Tutup',
		onClose
	}: Props = $props();

	function handleKeydown(event: KeyboardEvent): void {
		if (event.key === 'Escape' && open) {
			onClose();
		}
	}

	function getAccentClass(): string {
		if (type === 'success') {
			return 'bg-emerald-500';
		}

		if (type === 'error') {
			return 'bg-red-500';
		}

		return 'bg-blue-500';
	}

	function getBadgeClass(): string {
		if (type === 'success') {
			return 'bg-emerald-50 text-emerald-700 border-emerald-200';
		}

		if (type === 'error') {
			return 'bg-red-50 text-red-700 border-red-200';
		}

		return 'bg-blue-50 text-blue-700 border-blue-200';
	}

	function getLabel(): string {
		if (type === 'success') {
			return 'Berhasil';
		}

		if (type === 'error') {
			return 'Peringatan';
		}

		return 'Informasi';
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
	<div
		class="fixed inset-0 z-[100] flex items-center justify-center px-4 py-6"
		transition:fade={{
			duration: 150
		}}
	>
		<button
			type="button"
			aria-label="Tutup modal"
			class="absolute inset-0 bg-slate-950/50 backdrop-blur-[2px]"
			onclick={onClose}
		></button>

		<div
			role="dialog"
			aria-modal="true"
			aria-label={title}
			class="relative z-10 w-full max-w-lg overflow-hidden rounded-xl border border-slate-200 bg-white shadow-2xl"
			transition:scale={{
				duration: 170,
				start: 0.97
			}}
		>
			<div class={`h-1 w-full ${getAccentClass()}`}></div>

			<div
				class="flex items-start justify-between gap-4 border-b border-slate-200 px-5 py-4 sm:px-6"
			>
				<div>
					<div
						class={`mb-2 inline-flex rounded-md border px-2 py-1 text-[10px] font-semibold tracking-wide uppercase ${getBadgeClass()}`}
					>
						{getLabel()}
					</div>

					<h2 class="text-lg font-semibold text-slate-900">
						{title}
					</h2>
				</div>

				<button
					type="button"
					aria-label="Tutup"
					onclick={onClose}
					class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md text-xl leading-none text-slate-400 transition hover:bg-slate-100 hover:text-slate-700"
				>
					×
				</button>
			</div>

			<div class="px-5 py-5 sm:px-6">
				<p class="text-sm leading-6 text-slate-600">
					{message}
				</p>
			</div>

			<div class="flex justify-end gap-3 border-t border-slate-200 bg-slate-50 px-5 py-4 sm:px-6">
				<button
					type="button"
					onclick={onClose}
					class="inline-flex h-10 items-center justify-center rounded-md bg-blue-600 px-5 text-sm font-semibold text-white transition hover:bg-blue-700 focus:ring-4 focus:ring-blue-100 focus:outline-none"
				>
					{buttonText}
				</button>
			</div>
		</div>
	</div>
{/if}
