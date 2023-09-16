<script lang="ts">
	import { store, type Contributor } from "../../store";
	import { onMount } from "svelte";
	import AccessControl from "./access-control.svelte";
	import Unlink from "./unlink.svelte";

	export let repoName: string;
	const accessControl: string = "access-control";
	const unlinkRepo: string = "unlink-repo";

	let accessControlModal: HTMLDialogElement | undefined;
	let unlinkRepoModal: HTMLDialogElement | undefined;
	onMount(() => {
		accessControlModal = document.getElementById(
			accessControl,
		) as HTMLDialogElement;
		unlinkRepoModal = document.getElementById(
			unlinkRepo,
		) as HTMLDialogElement;
	});

	let contributors: Array<Contributor>;
	$: contributors = structuredClone($store.contributors);

	function handleOpenAccessControl() {
		contributors = structuredClone($store.contributors);
		document.querySelector("body")?.classList.add("overflow-y-hidden");
		accessControlModal?.showModal();
	}

	function handleOpenUnlinkRepo() {
		document.querySelector("body")?.classList.add("overflow-y-hidden");
		unlinkRepoModal?.showModal();
	}
</script>

<div class="dropdown dropdown-hover dropdown-end">
	<div class="btn btn-ghost bg-neutral hover:cursor-auto">
		<i class="fa-solid fa-gear"></i>
	</div>
	<ul
		class="dropdown-content menu bg-neutral rounded-box z-[1] w-52 p-2 shadow"
	>
		<li>
			<button on:click={handleOpenAccessControl}>
				<i class="fa-solid fa-user w-5"></i>
				Contributors
			</button>
		</li>
		<li>
			<a
				href={`https://github.com/${repoName}`}
				target="_blank"
				class=""
			>
				<i class="fa-brands fa-github w-5"></i>
				Git Repository
			</a>
		</li>
		{#if $store.is_owner}
			<button
				class="btn btn-outline"
				type="button"
				on:click={handleOpenUnlinkRepo}
			>
				<span>
					<i class="fa-solid fa-link-slash w-5"></i>
				</span>
				Unlink Repository
			</button>
		{/if}
	</ul>
</div>

<AccessControl
	modalID={accessControl}
	{repoName}
	{contributors}
/>

<Unlink
	id={unlinkRepo}
	{repoName}
/>
