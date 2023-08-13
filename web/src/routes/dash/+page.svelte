<script lang="ts">
	import { type User, AuthStore } from "@lib/auth";
	import { type Sort, fetchRepos } from "./requests";
	import Spinner from "@components/spinner.svelte";
	import Repo from "./components/Repo.svelte";
	import { onMount } from "svelte";

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
	onMount(() => {
		user = AuthStore.user();
	});

	let page: number = 0;
	let show: number = ShowDefault;
	let sort: Sort = SortDefault;
	const data = fetchRepos(page, sort, show);

	let search: string = "";
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

	<section class="">
		<h1 class="text-gradient text-2xl font-semibold">Repositories</h1>

		<div class="">
			<div class="grid grid-cols-2">
				<div>
					<label for="sort">Sort by: </label>
					<select name="sort" class="">
						{#each Object.entries(sortFunctions) as [name, sort]}
							<option value={sort}>
								{name}
							</option>
						{/each}
					</select>
				</div>

				<div>
					<label for="show">Show: </label>
					<select name="show" class="">
						{#each [10, 20, 30] as i}
							<option>
								{i}
							</option>
						{/each}
					</select>
				</div>
			</div>

			<div class="flex flex-col">
				<label for="repo-search">Search bar</label>
				<input name="repo-search" />
			</div>
		</div>

		<hr class="text-main border-main my-4 rounded-lg border" />

		{#await data}
			<div
				class="flex h-64 w-full flex-col items-center justify-center gap-3"
			>
				<Spinner class="stroke-main" />
				<p>Fetching data...</p>
			</div>
		{:then repos}
			{#if repos.length === 0}
				<p>You don't have any repository yet.</p>
			{:else}
				<ul>
					{#each repos.filter( (d) => d.full_name.includes(search ?? ""), ) as repo}
						<Repo  {repo} />
					{/each}
				</ul>
				<div class="">
					<button>prev</button>
					<p>{page}</p>
					<button>next</button>
				</div>
			{/if}
		{:catch error}
			<p>{error}</p>
		{/await}
	</section>
</main>
