<script lang="ts">
	import Modal from "@components/modal.svelte";
	import { store, type Contributor } from "../../store";
	import { afterUpdate, onMount } from "svelte";
	import cn from "classnames";

	export let repoName: string;

	let modal: HTMLDialogElement | undefined;
	onMount(() => {
		modal = document.getElementById("access-control") as HTMLDialogElement;
	});

	let ownerLogin: string;
	$: {
		for (const u of $store.contributors) {
			if (u.id === $store.owner_id) {
				ownerLogin = u.login;
				break;
			}
		}
	}

	let hasDiff: boolean = false;
	let contributors: Array<Contributor>;
	$: contributors = structuredClone($store.contributors);
	afterUpdate(() => {
		for (let i = 0; i < contributors.length; i++) {
			if (
				$store.contributors[i].write_access ===
				contributors[i].write_access
			)
				continue;
			hasDiff = true;
			return;
		}
		hasDiff = false;
	});

	let loading: boolean = false;

	// button handler
	function handleToggleWriteAccess(id: number) {
		return () => {
			contributors = contributors.map((c) => {
				if (c.id !== id) return c;
				c.write_access = !c.write_access;
				return c;
			});
		};
	}

	function handleOpen() {
		hasDiff = false;
		loading = false;
		contributors = structuredClone($store.contributors);
		document.querySelector("body")?.classList.add("overflow-y-hidden");
		modal?.showModal();
	}

	function handleClose() {
		document.querySelector("body")?.classList.remove("overflow-y-hidden");
		modal?.close();
	}

	function handleSubmit() {
		loading = true;
		console.log(contributors); // TODO: do something with this
		setTimeout(() => modal?.close(), 3000);
	}
</script>

<div class="dropdown dropdown-hover dropdown-end">
	<div class="btn btn-ghost bg-dark-100 hover:cursor-auto">
		<i class="fa-solid fa-gear"></i>
	</div>
	<ul
		class="dropdown-content menu bg-dark-100 rounded-box z-[1] w-52 p-2 shadow"
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

<Modal
	id="access-control"
	heading="Contributors / Access Control"
>
	<p class="mb-4">
		{#if $store.is_owner}
			All contributors will have
			<span class="font-semibold">read-only</span>
			access by default. To modify your contributor list, visit your
			<a
				href={`https://github.com/${repoName}/settings/access`}
				class="link font-semibold"
				target="_blank"
			>
				repository settings</a
			>.
		{:else}
			To request for access change, contact repository owner
			<a
				href={`https://github.com/${ownerLogin}`}
				class="link font-semibold"
				target="_blank">{ownerLogin}</a
			>.
		{/if}
	</p>

	<ul class="mb-10 ml-7 max-h-72 overflow-visible overflow-y-scroll">
		{#each $store.contributors as { id, login, avatar_url, write_access }, i (id)}
			<li class="my-5 flex items-center justify-start gap-3">
				<div class="w-8">
					<img
						src={avatar_url}
						alt={login}
						class="rounded-full"
					/>
				</div>

				<div>
					<div class="flex items-center gap-1">
						<p>
							{login}
						</p>
						{#if $store.is_owner}
							<input
								type="checkbox"
								class="toggle toggle-xs toggle-accent tooltip tooltip-right"
								data-tip="toggle write acccess"
								checked={contributors[i].write_access}
								on:click={handleToggleWriteAccess(id)}
							/>
						{/if}
					</div>

					<div class="flex items-start">
						{#if $store.owner_id === id}
							<p class="text-light/50 text-xs">owner /&nbsp;</p>
						{/if}
						{#if write_access}
							<p class="text-light/50 text-xs">read / write</p>
						{:else}
							<p class="text-light/50 text-xs">read-only</p>
						{/if}
					</div>
				</div>
			</li>
		{/each}
	</ul>

	<div class="flex justify-end gap-3">
		<button
			class="btn btn-ghost"
			on:click|preventDefault={handleClose}>Cancel</button
		>
		<button
			disabled={!hasDiff}
			type="button"
			on:click|preventDefault={handleSubmit}
			class={cn([
				"btn btn-primary w-24",
				{ "pointer-events-none": loading },
			])}
		>
			{#if loading}
				<span class="loading loading-spinner"></span>
			{:else}
				Commit
			{/if}
		</button>
	</div>
</Modal>
