<script lang="ts">
	import { store, type Contributor } from "../store";

	const SHOW: number = 5;

	let contributors: Array<Contributor>;
	$: contributors = $store.contributors.slice(0, SHOW);
</script>

<div class="avatar-group -space-x-6">
	{#each contributors as { id, name, avatar_url: src } (id)}
		<div class="avatar">
			<div
				class="tooltip w-12"
				data-tip={name}
			>
				<img
					{src}
					alt={name}
				/>
			</div>
		</div>
	{/each}

	{#if contributors.length > SHOW}
		<div class="avatar placeholder">
			<div class="bg-neutral-focus text-neutral-content w-12">
				<span>+{contributors.length - SHOW}</span>
			</div>
		</div>
	{/if}
</div>
