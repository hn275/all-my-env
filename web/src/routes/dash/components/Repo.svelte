<script lang="ts">
	import cx from "classnames";
	import type { Repository } from "../types";

	export let repo: Repository;

  function linkRepoURL(): string {
    const p: URLSearchParams  = new URLSearchParams({
      repo_id: repo.id.toString()
    });
    return "/repo/link?" + p.toString();
  }

  function repoURL(): string {
    return "/repo/" + String(repo.id);
  }

  const role: string = repo.is_owner ? "owner" : "collaborator";
  const githubHref: string = `https://www.github.com/${repo.full_name}`;

</script>

<li class="bg-dark-200 border border-dark-200 rounded-lg p-5 relative">
	<div class="flex gap-4 mb-6">
		<img
			src={repo.owner.avatar_url}
			alt={repo.owner.login}
			class="aspect-auto w-12 rounded-full"
		/>

    <div>
      <h3 class="text-lg block">{repo.name}</h3>
      <p class="text-light/70 text-sm inline">
        {role}
      </p>
    </div>
  </div>

  <div class="flex justify-between items-center">
    <div class="text-light/70 flex items-center gap-2">
      <a href={githubHref} target="_blank" class="mr-2">
        <i class="fa-brands fa-github fa-sm"></i>
      </a>

      {#if repo.fork}
        <i class="fa-solid fa-code-fork fa-xs"></i>
      {/if}

      {#if repo.linked}
        <i class="fa-solid fa-link fa-xs"></i>
      {/if}
    </div>


    {#if repo.linked}
      <div class="flex items-center justify-between py-2">
        <a href={repoURL()} class="text-sm text-main/80 hover:text-main transition">
          view 123 .env variables
        </a>
      </div>
    {:else}
      <a href={linkRepoURL()} class={cx([
        "bg-main/5 border-main", 
        "text-main rounded-lg border px-3 py-2 text-sm transition",
        "hover:bg-main/20 float-right hover:brightness-110"
      ])}>
        Connect Repository
      </a>
    {/if}
  </div>
</li>

