<script lang="ts">
	let {
		title,
		description
	}: {
		title: string;
		description?: string;
	} = $props();
</script>

<div class="bg-slate-800 px-4 pt-8 pb-24 text-white sm:px-6 lg:px-8">
	<div class="mx-auto max-w-7xl">
		<h1 class="text-xl font-semibold tracking-tight">
			{title}
		</h1>

		{#if description}
			<p class="mt-2 max-w-2xl text-sm leading-6 text-slate-300">
				{description}
			</p>
		{/if}
	</div>
</div>
