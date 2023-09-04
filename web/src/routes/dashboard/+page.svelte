<script lang="ts">
	import { fetchRepos } from "./requests";
	import type { Sort } from "./requests";
	import Spinner from "@components/spinner.svelte";
	import Repo from "./components/Repo.svelte";
	import Main from "@components/main.svelte";
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
	const Show: number = 30;
	let page: number = 1;

	let sort: Sort = SortDefault;
	async function handleSort(e: Event) {
		e.preventDefault();
		sort = (e.target as HTMLSelectElement)?.value as Sort;
		page = 1;
		await getRepos();
	}

	let repos: Array<Repository> = [];
	let loading: boolean = true;
	let error: string | undefined;
	onMount(async () => {
		page = 1;
		await getRepos();
	});

	let hasMoreRepo: boolean = true;
	async function getRepos() {
		try {
			loading = true;
			repos = await fetchRepos(page, sort, Show.toString());
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	let loadMoreLoading: boolean = false;
	onMount(() => {
		window.onscroll = async function (_: Event) {
			const bottom = window.innerHeight + Math.round(window.scrollY);
			const isBottom: boolean = bottom >= document.body.offsetHeight;
			if (!isBottom || !hasMoreRepo || loadMoreLoading) return;
			page++;
			try {
				loadMoreLoading = true;
				const res = await fetchRepos(page, sort, Show.toString());
				repos = [...repos, ...res];
				hasMoreRepo = res.length === Show;
				console.log(hasMoreRepo, page);
			} catch (e) {
				error = (e as Error).message;
			} finally {
				loadMoreLoading = false;
			}
		};
	});

	let search: string = "";
	function handleSearch(e: Event): void {
		search = (e.target as HTMLInputElement)?.value ?? "";
	}
</script>

<Main>
	<div class="mx-auto max-w-screen-2xl">
		<section class="h-full p-5 md:p-7">
			<h1 class="text-gradient mb-3 text-2xl font-semibold">
				Repositories
			</h1>
			<div class="mb-8">
				<div class="flex gap-5">
					<input
						type="text"
						placeholder="Search repositories"
						class="input input-bordered bg-neutral w-full flex-grow"
						on:input={handleSearch}
					/>
					<select
						class="select select-bordered bg-neutral text-light/70 font-normal"
						name="sort"
						on:change={handleSort}
					>
						{#each Object.entries(sortFunctions) as [name, s]}
							<option value={s} selected={sort === s}>
								{name}
							</option>
						{/each}
					</select>
				</div>
			</div>

			{#if error}
				<div class="text-dark-200 rounded-lg bg-red-400 p-5">
					<h2 class="inline text-lg font-bold">Whoops!</h2>
					<span>An error has occured:</span>
					<p>{error}</p>
				</div>
			{:else if loading}
				<div
					class="flex h-full min-h-[calc(100vh-420px)] w-full flex-col items-center justify-center gap-3"
				>
					<Spinner class="stroke-main" />
					<p>Fetching data...</p>
				</div>
			{:else if repos.length === 0}
				<p>You don't have any repository yet.</p>
			{:else}
				<ul
					class="grid grid-cols-1 gap-5 md:grid-cols-2 lg:grid-cols-3"
				>
					{#each repos.filter((r) =>
						r.full_name.includes(search ?? ""),
					) as repo (repo.id)}
						<Repo {repo} />
					{/each}
				</ul>
			{/if}

			{#if hasMoreRepo && loadMoreLoading && !loading}
				<div class="mb-4 mt-12 flex w-full justify-center">
					<div class="loading bg-main" />
				</div>
			{/if}
		</section>
	</div>
</Main>
