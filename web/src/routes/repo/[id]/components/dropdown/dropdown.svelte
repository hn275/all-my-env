<script lang="ts">
	import { store, type Contributor } from "../../store";
	import { onMount } from "svelte";
	import AccessControl from "./access-control.svelte";
	import UnlinkRepository from "./unlink-repository.svelte";

	export let repoName: string;

	const modalID: string = "access-control";
	let accessControlModal: HTMLDialogElement | undefined;

	const unlinkRepoID: string = "unlink-repo";
	let unlinkModal: HTMLDialogElement | undefined;
	onMount(() => {
		accessControlModal = document.getElementById(
			modalID,
		) as HTMLDialogElement;
		unlinkModal = document.getElementById(
			unlinkRepoID,
		) as HTMLDialogElement;
	});

	let contributors: Array<Contributor>;
	$: contributors = structuredClone($store.contributors);

	function handleOpenAccessControl() {
		contributors = structuredClone($store.contributors);
		document.querySelector("body")?.classList.add("overflow-y-hidden");
		accessControlModal?.showModal();
	}

	function handleOpenUnlink() {
		document.querySelector("body")?.classList.add("overflow-y-hidden");
		unlinkModal?.showModal();
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
			<li>
				<button on:click={handleOpenUnlink}>
					<i class="fa-solid fa-unlink w-5" />
					Unlink Repository
				</button>
			</li>
		{/if}
	</ul>
</div>

<AccessControl
	{modalID}
	{repoName}
	{contributors}
/>

<UnlinkRepository
	id={unlinkRepoID}
	{repoName}
/>
