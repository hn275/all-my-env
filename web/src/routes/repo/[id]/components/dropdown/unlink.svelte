<script lang="ts">
	import Modal from "@components/modal.svelte";
	import { store } from "../../store";
	import { makeUrl } from "@lib/url";
	import { apiFetch } from "@lib/requests";
	import { onMount } from "svelte";

	export let repoName: string;
	export let id: string;
	let confirmRepo: string;

	let modal: HTMLDialogElement | null;
	onMount(() => {
		modal = document.getElementById(id) as HTMLDialogElement;
	});

	let loading: boolean = false;
	let error: string | undefined;
	async function handleSubmit(e: Event) {
		e.preventDefault();
		try {
			loading = true;
			const url: string = makeUrl(`/repos/${$store.repoID}`);
			const r: RequestInit = { method: "DELETE" };
			const res = await apiFetch(url, r);
			if (res.status === 204) {
				window.location.replace("/dashboard");
				return;
			}
			const payload: EnvHub.Error = await res.json();
			error = payload.message;
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}
</script>

<Modal
	{id}
	heading="Unlink Repository"
>
	<p class="mb-4">
		You're about to delete the repository
		<span class="font-bold">
			{repoName}
		</span>
		from your EnvHub account. You may reconnect this repository, but the saved
		variables will be
		<span class="font-bold"> permanently deleted. </span>
		<br />
	</p>
	<div class="alert alert-error mb-5">
		<i class="fa-solid fa-triangle-exclamation" />
		<p class="">
			<span class="font-bold">Warning:</span> this action is not reversible.
		</p>
	</div>

	<p class="mb-2">
		Type
		<span class="font-bold">`{repoName}`</span>
		to confirm:
	</p>

	<div class="flex h-[70px] flex-col items-start justify-start">
		<input
			type="text"
			class="input input-bordered w-full"
			bind:value={confirmRepo}
			placeholder={repoName}
			required
		/>
		{#if error}
			<p class="text-error text-xs">{error}</p>
		{/if}
	</div>
	<div class="modal-action">
		<button
			class="btn btn-ghost"
			type="button"
			on:click={() => modal?.close()}
		>
			Cancel
		</button>
		<button
			class="btn btn-error w-44"
			type="submit"
			disabled={repoName !== confirmRepo}
			on:click|preventDefault={handleSubmit}
		>
			{#if loading}
				<span class="loading" />
			{:else}
				Delete my variable
			{/if}
		</button>
	</div>
</Modal>
