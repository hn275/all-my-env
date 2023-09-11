<script lang="ts">
	import Main from "@components/main.svelte";
	import type { Breadcrumbs } from "@lib/types";
	import { onMount } from "svelte";
	import type { Route } from "./+page.server";
	import cx from "classnames";
	import { Table, NewModal, Contributors, Dropdown } from "./components";
	import { store } from "./store";

	export let data: Route;
	let breadcrumbs: Array<Breadcrumbs> | undefined;
	let repoName: string;
	onMount(() => {
		const url: URL = new URL(window.location.href);
		repoName = url.searchParams.get("name")!;
		repoName = decodeURIComponent(repoName);
		breadcrumbs = [{ text: repoName, href: url.toString() }];
	});
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
				{#if $store.write_access}
					<NewModal />
				{/if}
			</div>
		</div>
	</section>

	<section class="mx-auto max-w-screen-2xl">
		<div
			class="bg-neutral card mt-5 min-h-[400px] w-full overflow-x-auto p-5 shadow-xl"
		>
			<Table repoID={data.id} />
		</div>
	</section>
</Main>

<style lang="postcss">
</style>
