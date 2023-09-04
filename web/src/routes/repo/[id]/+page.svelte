<script lang="ts">
	import Main from "@components/main.svelte";
	import type { Breadcrumbs } from "@lib/types";
	import { onMount } from "svelte";
	import { getVariables } from "./services";
	import type { Route } from "./+page.server";
	import cx from "classnames";
	import {
		DeleteRepo,
		Table,
		NewModal,
		Contributors,
		Dropdown,
	} from "./components";
	import type { RepositoryEnv } from "./store";
	import { store } from "./store";

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

	let state: RepositoryEnv;
	$: state = $store;
</script>

<Main {breadcrumbs}>
	<section>
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
				<Contributors />
				<Dropdown {repoName} />
				{#if state.is_owner}
					<DeleteRepo {repoName} />
				{/if}
				{#if state.write_access}
					<NewModal />
				{/if}
			</div>
		</div>
	</section>

	<section class="mx-auto max-w-screen-2xl">
		<div
			class="bg-neutral card mt-5 min-h-[400px] w-full overflow-x-auto p-5 shadow-xl"
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
</Main>

<style lang="postcss">
</style>
