<script lang="ts">
	import Variable from "./variable.svelte";
	import { store } from "../store";
	import Row from "./row.svelte";
	import AddButton from "./new-modal.svelte";
</script>

<Row>
	<div />
	<h3>Key</h3>
	<h3>Value</h3>
	<h3>Created At</h3>
	<h3>Last Modified</h3>
</Row>
{#each $store.variables as variable, i (variable.id)}
	<Variable {i} {...variable} />
{/each}
{#if $store.variables.length === 0}
	<div
		class="flex h-full min-h-[400px] w-full flex-col items-center justify-center gap-3"
	>
		<p class="text-light/50">No variables stored</p>
		<AddButton />
	</div>
{/if}

<style lang="postcss">
	h3 {
		@apply font-semibold text-light/70;
		@apply ml-2;
	}
</style>
