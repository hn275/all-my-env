<script lang="ts">
  import Main from "@components/main.svelte";
  import type { Breadcrumbs } from "@lib/types";
	import { onMount } from "svelte";
  import { getVariables } from "./requests";
  import { store, type Variable } from "./store"
	import cx from "classnames";

  let breadcrumbs: Array<Breadcrumbs> | undefined;
  let repoID: number;
  let rsp: Promise<void>;
  let repoName: string;
  onMount(async () => {
    // breadcrumbs
    const url: URL = new URL(window.location.href);
    repoName = url.searchParams.get("name")!;
    repoName = decodeURIComponent(repoName);
    breadcrumbs = [{ text: repoName, href: url.toString() }];

    // fetch variables
    const s = url.searchParams.get("id") ?? "";
    repoID = parseInt(s);
    rsp = getVariables(repoID);
  })

  let variables: Array<Variable> = [];
  store.subscribe((v) => variables = v);
</script>

<Main {breadcrumbs}>
  <section class="bg-dark-200">
    <div class={cx([
      "flex items-center justify-between",
      "mx-auto max-w-screen-2xl px-4 py-5 text-sm",
    ])}>
      <h1 class="text-gradient text-2xl font-semibold">
        {repoName}
      </h1>

      <div class="flex gap-5">
        <a
          href={`https://github.com/${repoName}`} 
          target="_blank"
          class="btn btn-outline border-dark-100 text-xs"
        >
          Git Repository
        </a>
        <button class="btn btn-primary text-xs">
          <i class="fa-solid fa-plus"></i>
          Add
        </button>
      </div>
    </div>
  </section>

  <section class="mx-auto max-w-screen-2xl">
    {#await rsp}
      <p>Loading</p>
      {:then}
      {#if variables.length === 0}
        alskdfj
      {:else}
        <ul>
          {#each variables as variable (variable.id)}
            <li></li>
          {/each}
        </ul>
      {/if}
      {:catch e}
      <p>{e.message}</p>
    {/await}
  </section>
</Main>
