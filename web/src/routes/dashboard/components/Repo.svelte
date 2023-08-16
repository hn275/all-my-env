<script lang="ts">
	import cx from "classnames";
	import type { Repository } from "../types";
	import Icon from "./RepoIcon.svelte";
	import { linkRepo } from "../requests";

	export let repo: Repository;

	const param = new URLSearchParams({
		name: encodeURIComponent(repo.full_name),
	});
	const repoHref: string = `/repo/${repo.id}?${param.toString()}`;

	const role: string = repo.is_owner ? "owner" : "collaborator";
	const githubHref: string = `https://www.github.com/${repo.full_name}`;

	let loading: boolean = false;
	let error: string | undefined;
	async function handleLinkRepo(): Promise<void> {
		try {
			loading = true;
			await linkRepo(repo);
			repo.linked = true;
		} catch (err) {
			error = (err as Error).message;
		} finally {
			loading = false;
			if (error) {
				setTimeout(() => {
					error = undefined;
				}, 3000);
			}
		}
	}
</script>

<li
	class={cx([
		"bg-dark-200  border-light/10 relative rounded-lg border p-5 shadow-lg",
		"flex flex-col justify-between",
	])}
>
	<div class="mb-6 flex gap-4">
		<img
			src={repo.owner.avatar_url}
			alt={repo.owner.login}
			class="aspect-auto w-12 rounded-full"
			loading="lazy"
		/>

		<div>
			<h3 class="block w-[17ch] truncate text-ellipsis md:w-[25ch] lg:w-[45ch]">
				{repo.name}
			</h3>
			<p class="text-light/70 inline text-sm">
				{role}
			</p>
		</div>
	</div>

	<div class="flex items-center justify-between">
		<div class="text-light/70 flex items-center gap-2">
			<Icon show={true} tooltip={repo.full_name}>
				<a href={githubHref} target="_blank" class="mr-2">
					<i class="fa-brands fa-github fa-sm" />
				</a>
			</Icon>

			<Icon show={repo.fork} tooltip="forked repo">
				<i class="fa-solid fa-code-fork fa-xs" />
			</Icon>

			<Icon show={repo.linked} tooltip="repo connected">
				<i class="fa-solid fa-link fa-xs" />
			</Icon>

			<Icon show={repo.private} tooltip="private">
				<i class="fa-solid fa-lock fa-xs" />
			</Icon>
		</div>

		{#if repo.linked}
			<div class="flex items-center justify-between py-2">
				<a
					href={repoHref}
					class={cx([
						"link link-primary text-xs",
						"bg-transparent text-sm transition hover:bg-transparent",
					])}
				>
					see {repo.variable_counter} variable(s)
				</a>
			</div>
		{:else if repo.is_owner}
			<button
				on:click={handleLinkRepo}
				class={cx([
					"text-main border-main bg-main/5 hover:bg-main/10 rounded-lg border p-2 text-sm transition",
					{ "loading loading-dots": loading },
				])}
			>
				{#if !loading}
					Connect repo
				{/if}
			</button>
		{:else}
			<p class="text-light/50 text-xs">not connected</p>
		{/if}
	</div>
</li>

{#if error}
	<div class="toast toast-bottom toast-start z-[1000]">
		<div class="alert alert-error flex flex-col items-start gap-1">
			<h6 class="font-semibold">Whoops!</h6>
			<p>{error}</p>
		</div>
	</div>
{/if}
