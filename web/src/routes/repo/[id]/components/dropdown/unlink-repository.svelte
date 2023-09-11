<script lang="ts">
	import { afterUpdate, onMount } from "svelte";
	import type { RepositoryEnv } from "../../store";
	import { store } from "../../store";
	import { apiFetch } from "@lib/requests";
	import { makeUrl } from "@lib/url";
	import Modal from "@components/modal.svelte";

	export let repoName: string;
	export let id: string;

	let repo: RepositoryEnv;
	onMount(() => {
		const unsub = store.subscribe((s) => (repo = s));
		return unsub;
	});

	let loading: boolean = false;
	let error: string | undefined;

	async function handleSubmit(e: Event) {
		e.preventDefault();
		if (!repo.repoID) {
			return;
		}

		if (confirmRepo !== repoName) {
			error = "Please confirm the repository Name.";
		}

		try {
			loading = true;
			const url: string = makeUrl(`/repos/${repo.repoID}/unlink`, {});
			const res = await apiFetch(url, {
				method: "DELETE",
			});

			if (res.status === 204) {
				window.location.replace("/dashboard");
			} else {
				const responseBody = await res.json();
				error = responseBody.message;
			}
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	afterUpdate(() => {
		if (!error) return;
		setTimeout(() => {
			error = undefined;
		}, 3000);
	});

	let confirmRepo: string = "";
</script>

<Modal
	heading="Unlink Repository"
	{id}
>
	<h3 class="mb-7 text-lg font-bold">Delete repository</h3>
	<p class="mb-4">
		You're about to unlink the repository
		<span class="font-bold">
			{repoName}
		</span>
		from your EnvHub account. You may link this repository at anytime in the
		future, but the saved variables will be
		<span class="font-bold">permanently deleted.</span>
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
		<button class="btn btn-ghost"> Cancel </button>
		<button
			class="btn btn-error w-44"
			type="submit"
			disabled={repoName !== confirmRepo}
			on:click|preventDefault={handleSubmit}
		>
			{#if loading}
				<span class="loading" />
			{:else}
				Delete My Repo
			{/if}
		</button>
	</div>
</Modal>
