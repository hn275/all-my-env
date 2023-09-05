<script lang="ts">
	import { store, type Contributor } from "../../store";
	import { onMount } from "svelte";
	import AccessControl from "./access-control.svelte";

	export let repoName: string;
	const modalID: string = "access-control";

	let modal: HTMLDialogElement | undefined;
	onMount(() => {
		modal = document.getElementById(modalID) as HTMLDialogElement;
	});

	let contributors: Array<Contributor>;
	$: contributors = structuredClone($store.contributors);

	function handleOpen() {
		contributors = structuredClone($store.contributors);
		document.querySelector("body")?.classList.add("overflow-y-hidden");
		modal?.showModal();
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
			<button on:click={handleOpen}>
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
	</ul>
</div>

<AccessControl
	{modalID}
	{repoName}
	{contributors}
/>
