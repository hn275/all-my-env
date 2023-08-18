<script lang="ts">
	import Main from "@components/main.svelte";
	import type { Breadcrumbs } from "@lib/types";
	import { onMount } from "svelte";
	import { getVariables } from "./requests";
	import { store } from "./store";
	import type { Route } from "./+page.server";
	import cx from "classnames";
	import AddButton from "./addBtn.svelte";

	export let data: Route;
	let breadcrumbs: Array<Breadcrumbs> | undefined;
	let rsp: Promise<void>;
	let repoName: string;
	onMount(async () => {
		// breadcrumbs
		const url: URL = new URL(window.location.href);
		repoName = url.searchParams.get("name")!;
		repoName = decodeURIComponent(repoName);
		breadcrumbs = [{ text: repoName, href: url.toString() }];

		// fetch variables
		rsp = getVariables(data.id);
	});
</script>

<Main {breadcrumbs}>
	<section class="bg-dark-200">
		<div
			class={cx([
				"flex items-center justify-between",
				"mx-auto max-w-screen-2xl px-4 py-5 text-sm",
			])}
		>
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
				<AddButton repoID={data.id} writeAccess={$store.write_access} />
			</div>
		</div>
	</section>

	<section class="m-4 mx-auto max-w-screen-2xl p-5 md:m-8">
		<h2 class="text-xl font-bold">Environment Variables</h2>
		<div
			class={cx([
				"bg-dark-200 mt-5 min-h-[400px]",
				"w-full overflow-x-auto rounded-md p-5",
			])}
		>
			<table class="table-xs md:table">
				<thead class="text-xs md:text-sm">
					<tr>
						<th />
						<th class="mr-8">Key</th>
						<th>Value</th>
						<th>Created By</th>
						<th>Last Modified</th>
					</tr>
				</thead>
				{#await rsp then}
					<tbody>
						{#each $store.variables as v, i (v.id)}
							<tr>
								<td>{i + 1}</td>
								<td>{v.key}</td>
								<td>{v.value}</td>
								<td>{v.created_at}</td>
								<td>{v.updated_at}</td>
							</tr>
						{/each}
					</tbody>
				{/await}
			</table>
			{#await rsp}
				<div
					class="flex h-full min-h-[400px] w-full flex-col items-center justify-center gap-5"
				>
					<span class="loading loading-lg text-main" />
					<p>Getting variables...</p>
				</div>
			{:then}
				{#if $store.variables.length === 0}
					<div
						class="flex h-full min-h-[400px] w-full flex-col items-center justify-center gap-3"
					>
						<p class="text-light/50">No variables stored</p>
						<AddButton repoID={data.id} writeAccess={$store.write_access} />
					</div>
				{/if}
			{:catch e}
				<div class="flex h-52 w-full flex-col items-center justify-center">
					<p class="text-error">{e.message}</p>
				</div>
			{/await}
		</div>
	</section>
</Main>

<style lang="postcss">
</style>
