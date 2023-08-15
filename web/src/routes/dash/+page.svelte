<script lang="ts">
	import { type User, AuthStore } from "@lib/auth";
	import { type Sort, fetchRepos } from "./requests";
	import Spinner from "@components/spinner.svelte";
	import Repo from "./components/Repo.svelte";
  import { onMount } from "svelte";
	import type { Repository } from "./types";

	// sort
	const SortDefault: Sort = "pushed";
	const sortFunctions: Record<string, Sort> = {
		Created: "created",
		Pushed: "pushed",
		Updated: "updated",
		Name: "full_name",
	};

	// page limit
	const ShowDefault: number = 30;

	let user: User | undefined;
	let page: number = 1;
	let show: number = ShowDefault;

	let sort: Sort = SortDefault;
  async function handleSort(e: Event) {
    e.preventDefault();
    sort = (e.target as HTMLSelectElement)?.value as Sort;
    await getRepos();
  }

  let repos: Array<Repository> = [];
  let loading: boolean = true;
  let error: string | undefined;
	onMount(async () => {
		user = AuthStore.user();
    if (!user) {
      AuthStore.refreshSession();
      return
    }
    await getRepos();
	});

  let hasMoreRepo: boolean = true;
  async function getRepos() {
    try {
      loading = true;
      repos = await fetchRepos(page, sort, show.toString());
    } catch (e) {
      error = (e as Error).message;
    } finally {
      loading = false;
    }
  }

  let loadMoreLoading: boolean = false;
  async function loadMore() {
    if (!hasMoreRepo) return;
    page++;
    try {
      loadMoreLoading = true;
      const res = await fetchRepos(page, sort, show.toString());
      repos = [...repos, ...res];
      hasMoreRepo = res.length === show;
    } catch (e) {
      error = (e as Error).message;
    } finally {
      loadMoreLoading = false;
    }
  }

  onMount(() => {
    window.onscroll = async function(_: Event) {
      const bottom: boolean = (window.innerHeight + Math.round(window.scrollY)) >= document.body.offsetHeight
      if (bottom) {
        await loadMore();
      }
    };
  })


	let search: string = "";
  function handleSearch(e: Event) {
    search = (e.target as HTMLInputElement)?.value ?? "";
  }

</script>

<main class="px-4 md:px-10">
	<section class="mt-5 px-4">
		<div class="flex flex-col items-center justify-center gap-2 md:w-max">
			<img
				alt="hn275"
				role="presentation"
				src={user?.avatar_url}
				class="w-20 rounded-full"
			/>
			<div>
				<h2 class="inline-block text-lg font-semibold">{user?.name}</h2>
				&nbsp;
				<span class="text-xs">({user?.login})</span>
			</div>
		</div>
	</section>

	<section class="bg-dark-100 mt-12 rounded-lg p-5 md:p-7">
		<h1 class="text-gradient mb-3 text-2xl font-semibold">Repositories</h1>

    <div class="">
      <div>
        <label for="sort">Sort by: </label>
        <select name="sort" class="" on:change={handleSort}>
          {#each Object.entries(sortFunctions) as [name, s]}
            <option value={s} selected={sort === s}>
              {name}
            </option>
          {/each}
        </select>
      </div>

			<div class="flex flex-col">
				<label for="repo-search">Search bar</label>
        <input name="repo-search" on:input={handleSearch}/>
			</div>
		</div>

		<hr class="text-main border-main my-6 rounded-lg border" />

    {#if error}
      <div class="text-dark-200 rounded-lg bg-red-400 p-5">
        <h2 class="inline text-lg font-bold">Whoops!</h2>
        <span>An error has occured:</span>
        <p>{error}</p>
      </div>

    {:else if loading}
      <div
        class="flex h-64 w-full flex-col items-center justify-center gap-3"
      >
        <Spinner class="stroke-main" />
        <p>Fetching data...</p>
      </div>

    {:else}
      {#if repos.length === 0}
        <p>You don't have any repository yet.</p>

      {:else}
        <ul class="grid grid-cols-1 gap-5 md:grid-cols-2 lg:grid-cols-3">
          {#each repos.filter((d) => d.full_name.includes(search ?? "")) as repo (repo.id)}
            <Repo {repo} />
          {/each}
        </ul>
      {/if}
    {/if}

    {#if hasMoreRepo && loadMoreLoading && !loading}
      <div class="mt-5 flex w-full justify-center">
        <div class="loading bg-main"></div>
      </div>
    {/if}
  </section>
</main>
