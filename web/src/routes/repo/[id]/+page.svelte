<script lang="ts">
	import Main from "@components/main.svelte";
	import type { Breadcrumbs } from "@lib/types";
	import { onMount } from "svelte";
	import { getVariables } from "./services";
	import type { Route } from "./+page.server";
	import cx from "classnames";
	import { DeleteModal, Table, NewModal, UnlinkRepo } from "./components";

	export let data: Route;
	let breadcrumbs: Array<Breadcrumbs> | undefined;
	let rsp: Promise<void>;
	let repoName: string;
	onMount(async () => {
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
			<div class="flex gap-3">
				<UnlinkRepo />
				<a
					href={`https://github.com/${repoName}`}
					target="_blank"
					class="btn btn-outline border-dark-100 text-xs"
				>
					Git Repository
				</a>
				<button
					class={cx([
						"btn btn-outline bg-main/0 hover:bg-main/20 border-main/10",
						"text-main/80 hover:text-main text-xs",
					])}
				>
					Download file
				</button>
				<NewModal />
			</div>
		</div>
	</section>

	<section class="m-4 mx-auto max-w-screen-2xl p-5">
		<h2 class="text-xl font-bold">Environment Variables</h2>
		<div
			class={cx([
				"bg-dark-200 mt-5 min-h-[400px] shadow shadow-black/20",
				"border-light/20 border",
				"w-full overflow-x-auto rounded-md p-5",
			])}
		>
			{#await rsp}
				<div
					class={cx([
						"flex h-full min-h-[400px] w-full flex-col",
						"items-center justify-center gap-5",
					])}
				>
					<span class="loading loading-lg text-main" />
					<p>Getting variables...</p>
				</div>
			{:then}
				<Table />
			{:catch e}
				<div
					class="flex h-52 w-full flex-col items-center justify-center"
				>
					<p class="text-error">{e.message}</p>
				</div>
			{/await}
		</div>
	</section>

	<DeleteModal />
</Main>

<style lang="postcss">
</style>
