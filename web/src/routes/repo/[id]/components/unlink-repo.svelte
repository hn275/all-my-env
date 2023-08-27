<script lang="ts">
	import cx from "classnames";
	import { afterUpdate, onMount } from "svelte";
	import type { RepositoryEnv } from "../store";
	import { store } from "../store";
	import { apiFetch } from "@lib/requests";
	import { makeUrl } from "@lib/url";

	export let repoName: string;

	let repo: RepositoryEnv;
	onMount(() => {
		const unsub = store.subscribe((s) => (repo = s));
		return unsub;
	});

	let modal: HTMLDialogElement | null;
	onMount(() => {
		modal = document.getElementById("del-repo") as HTMLDialogElement;
	});

	let loading: boolean = false;
	let error: string | undefined;
	async function handleSubmit(e: Event) {
		if (!repo.repoID) throw new Error("repository not found.");
		e.preventDefault();
		try {
			loading = true;
			const url: string = makeUrl(`/repos/${repo.repoID}/delete`);
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

	afterUpdate(() => {
		if (!error) return;
		setTimeout(() => {
			error = undefined;
		}, 3000);
	});

	let confirmRepo: string = "";
</script>

<button
	class={cx([
		"btn btn-outline border-error/20 text-error/50 hover:text-error",
		"hover:border-error/80 hover:bg-dark-200 disabled:bg-dark-200 disabled:text-error/30 text-xs",
	])}
	disabled={!$store.write_access}
	on:click={() => modal?.showModal()}
>
	<i class="fa-solid fa-link-slash" />
	Delete repository
</button>

<dialog
	id="del-repo"
	class="modal"
>
	<form
		method="dialog"
		class="modal-box bg-dark-200"
	>
		<div class="bg-dark-200">
			<h3 class="mb-7 text-lg font-bold">Delete repository</h3>
			<p class="mb-4">
				You're about to delete the repository
				<span class="font-bold">
					{repoName}
				</span>
				from your EnvHub account. You may reconnect this repository, but
				the saved variables will be
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
		</div>
	</form>
	<form
		id="new-var-close"
		method="dialog"
		class="modal-backdrop"
	>
		<button id="new-var-close" />
	</form>
</dialog>
