<script lang="ts">
	import { store, type Contributor } from "../store";

	const SHOW: number = 4;

	let contributors: Array<Contributor>;
	$: contributors = $store.contributors;

	let contributorCounter: string;
	$: contributorCounter = `+${contributors.length - SHOW}`;

	let ownerID: number;
	$: ownerID = $store.owner_id;
</script>

<ul
	class="hidden md:flex avatar-group -space-x-5 transition-all hover:-space-x-2"
>
	{#each contributors.slice(0, SHOW) as { id, login, avatar_url: src } (id)}
		<li class="avatar transition-all duration-300">
			<a
				href={`https://github.com/${login}`}
				target="_blank"
				class="w-10"
			>
				<img
					{src}
					alt={login}
				/>
			</a>
		</li>
	{/each}
	{#if contributors.length > SHOW}
		<li
			class="avatar placeholder bg-neutral group-hover:border-primary transition-all"
		>
			<div class="h10 w-10 text-neutral-content text-xs font-bold">
				{contributorCounter}
			</div>
		</li>
	{/if}
</ul>

<style>
</style>
