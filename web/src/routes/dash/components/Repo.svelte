<script lang="ts">
	import type { Repository } from "../types";

	export let repo: Repository;

	const repoParam: string = new URLSearchParams({
		repo_name: encodeURI(repo.full_name),
	});

</script>

<li class="bg-dark-100">
	<div>
		<img
			src={repo.owner.avatar_url}
			alt={repo.owner.login}
			class="aspect-auto w-8 rounded-full"
		/>

		<a href="/repoid">{repo.name}</a>
		{#if !repo.is_owner}
			<p>Collaborator</p>
		{:else}
			{#if repo.linked}
				<p>Linked</p>
      {:else}
        <a href="/link">Link repo</a>
      {/if}
		{/if}

		{#if repo.fork}
			<p>forked</p>
		{/if}
	</div>

	<a href={`https://www.github.com/${repo.full_name}`} target="_blank">
		GitHub Logo
	</a>
</li>
