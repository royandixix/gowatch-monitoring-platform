<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';

	type InputType = HTMLInputAttributes['type'];

	type InputAutocomplete = HTMLInputAttributes['autocomplete'];

	interface Props {
		id: string;
		label: string;
		type?: InputType;
		placeholder?: string;
		value?: string;
		required?: boolean;
		disabled?: boolean;
		error?: string;
		hint?: string;
		autocomplete?: InputAutocomplete;
	}

	let {
		id,
		label,
		type = 'text',
		placeholder = '',
		value = $bindable(''),
		required = false,
		disabled = false,
		error = '',
		hint = '',
		autocomplete
	}: Props = $props();
</script>

<div>
	<label for={id} class="mb-2 block text-xs font-semibold text-slate-600">
		{label}

		{#if required}
			<span class="text-red-500"> * </span>
		{/if}
	</label>

	<input
		{id}
		{type}
		{placeholder}
		{required}
		{disabled}
		{autocomplete}
		bind:value
		aria-invalid={error ? 'true' : undefined}
		aria-describedby={error || hint ? `${id}-message` : undefined}
		class={`h-11 w-full rounded-md border bg-white px-3 text-sm text-slate-700 transition outline-none placeholder:text-slate-400 disabled:cursor-not-allowed disabled:bg-slate-50 disabled:text-slate-400 ${
			error
				? 'border-red-300 focus:border-red-400 focus:ring-4 focus:ring-red-50'
				: 'border-slate-200 focus:border-sky-400 focus:ring-4 focus:ring-sky-50'
		}`}
	/>

	{#if error}
		<p id={`${id}-message`} class="mt-1.5 text-xs text-red-500">
			{error}
		</p>
	{:else if hint}
		<p id={`${id}-message`} class="mt-1.5 text-xs text-slate-400">
			{hint}
		</p>
	{/if}
</div>
