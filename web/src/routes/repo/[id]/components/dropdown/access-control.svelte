<script lang="ts">
	import Modal from "@components/modal.svelte";
	import { store, type Contributor } from "../../store";
	import { afterUpdate, onMount } from "svelte";
	import cn from "classnames";
	import { handlePermission } from "../../services";

	export let repoName: string;
	export let modalID: string;
	export let contributors: Contributor[];
	let ownerLogin: string;
	$: {
		for (const u of $store.contributors) {
			if (u.id === $store.owner_id) {
				ownerLogin = u.login;
				break;
			}
		}
	}

	let loading: boolean = false;
	let modal: HTMLDialogElement | undefined;
	onMount(() => {
		modal = document.getElementById(modalID) as HTMLDialogElement;
	});

	let hasDiff: boolean = false;
	afterUpdate(() => {
		for (let i = 0; i < contributors.length; i++) {
			if (
				contributors[i].write_access ===
				$store.contributors[i].write_access
			) {
				continue;
			} else {
				hasDiff = true;
				return;
			}
		}
		hasDiff = false;
	});

	let diffIndex: boolean[] = new Array(contributors.length).fill(false);
	$: {
		const c = $store.contributors;
		for (let i = 0; i < c.length; i++) {
			diffIndex[i] = contributors[i].write_access !== c[i].write_access;
			diffIndex = diffIndex;
		}
	}

	// button handler
	function handleToggleWriteAccess(id: number) {
		return () => {
			let i: number = 0;
			for (const c of contributors) {
				if (c.id === id) {
					break;
				}
				i++;
			}
			contributors[i].write_access = !contributors[i].write_access;
		};
	}

	async function handleSubmit() {
		loading = true;
		const userIDs: number[] = [];
		for (let i = 0; i < contributors.length; i++) {
			if (contributors[i].write_access) userIDs.push(contributors[i].id);
		}

		try {
			await handlePermission($store.repoID!, userIDs);
			const html = document.querySelector("body");
			html?.classList.remove("overflow-y-hidden");
			modal?.close();
		} catch (e) {
			console.error(e);
		} finally {
			loading = false;
		}
	}

	function handleClose() {
		document.querySelector("body")?.classList.remove("overflow-y-hidden");
		modal?.close();
	}
</script>

<Modal
	id={modalID}
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
		{#each contributors as { id, login, avatar_url, write_access }, i (id)}
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
						{#if $store.is_owner && $store.owner_id !== id}
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
						{#if diffIndex[i]}
							<span class="text-light/50 ml-2 text-xs">[+]</span>
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
