<script lang="ts">
	import { fetchRepos } from "./requests";
	import type { Sort } from "./requests";
	import Repo from "./components/Repo.svelte";
	import Main from "@components/main.svelte";
	import { onMount } from "svelte";
	import type { Repository } from "./types";
	import { Auth, oauth } from "@lib/auth";
	import { refresh } from "../index/lib/auth";

	// sort
	const SortDefault: Sort = "pushed";
	const sortFunctions: Record<string, Sort> = {
		Created: "created",
		Pushed: "pushed",
		Updated: "updated",
		Name: "full_name",
	};

	let rsp: Promise<void>;

	// page limit
	const Show: number = 30;
	let page: number = 1;

	let sort: Sort = SortDefault;
	function handleSort(e: Event) {
		e.preventDefault();
		sort = (e.target as HTMLSelectElement)?.value as Sort;
		page = 1;
		rsp = getRepos();
	}

	let repos: Array<Repository> = [];
	onMount(async () => {
		let user = Auth.user();
		if (!user) {
			oauth();
			return;
		}

		if (!Auth.isRefreshed()) {
			user = await refresh(user.access_token);
			if (!user) {
				oauth();
				return;
			}

			Auth.login(user);
			Auth.refresh();
		}
		page = 1;
		rsp = getRepos();
	});

	let hasMoreRepo: boolean = true;
	async function getRepos() {
		repos = await fetchRepos(page, sort, Show.toString());
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
			} catch (e) {
				console.error(e);
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
							<option
								value={s}
								selected={sort === s}
							>
								{name}
							</option>
						{/each}
					</select>
				</div>
			</div>

			{#await rsp}
				<div
					class="flex h-full min-h-[calc(100vh-420px)] w-full flex-col items-center justify-center gap-3"
				>
					<span
						class="loading loading-lg loading-ring text-primary"
					/>
				</div>
			{:then}
				<ul
					class="grid grid-cols-1 gap-5 md:grid-cols-2 lg:grid-cols-3"
				>
					{#each repos.filter( (r) => r.full_name.includes(search ?? ""), ) as repo (repo.id)}
						<Repo {repo} />
					{/each}
				</ul>
			{:catch e}
				<div class="text-error-content bg-error rounded-lg p-5">
					<h2 class="inline text-lg font-bold">Whoops!</h2>
					<span>An error has occured:</span>
					<p>{e}</p>
				</div>
			{/await}

			{#if hasMoreRepo && loadMoreLoading}
				<div class="mb-4 mt-12 flex w-full justify-center">
					<div class="loading loading-ring text-primary" />
				</div>
			{/if}
		</section>
	</div>
</Main>
