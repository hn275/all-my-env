<script lang="ts">
	import NewVariable from "./new-var.svelte";
	import Row from "./table-row.svelte";
	import { store } from "./store";

	export let repoID: number;
</script>

<table class="table-xs md:table mx-auto">
	<thead class="text-xs md:text-sm">
		<tr>
			<th />
			<th class="w-7">Key</th>
			<th>Value</th>
			<th>Created At</th>
			<th>Last Modified</th>
		</tr>
	</thead>
	<tbody>
		{#each $store.variables as v, i (v.id)}
			<Row {i} {...v} />
		{/each}
	</tbody>
</table>
{#if $store.variables.length === 0}
	<div
		class="flex h-full min-h-[400px] w-full flex-col items-center justify-center gap-3"
	>
		<p class="text-light/50">No variables stored</p>
		<NewVariable {repoID} writeAccess={$store.write_access} />
	</div>
{/if}
