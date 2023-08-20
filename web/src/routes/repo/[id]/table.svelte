<script lang="ts">
	import NewVariable from "./new-var.svelte";
	import Variable from "./variable.svelte";
	import { store } from "./store";
	import Row from "./row.svelte";

	export let repoID: number;
</script>

<Row>
	<div />
	<h3>Key</h3>
	<h3>Value</h3>
	<h3>Created At</h3>
	<h3>Last Modified</h3>
</Row>
{#each $store.variables as variable, i (variable.id)}
	<Row>
		<Variable {i} {...variable} />
	</Row>
{/each}
{#if $store.variables.length === 0}
	<div
		class="flex h-full min-h-[400px] w-full flex-col items-center justify-center gap-3"
	>
		<p class="text-light/50">No variables stored</p>
		<NewVariable {repoID} writeAccess={$store.write_access} />
	</div>
{/if}

<style lang="postcss">
	h3 {
		@apply font-semibold text-light/70;
		@apply ml-2;
	}
</style>
