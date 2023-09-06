<script lang="ts">
	import { handleDelete } from "../../services";
	import { onMount } from "svelte";

	export let variableKey: string;
	export let variableID: string;
	export let repoID: number | undefined;
	export let _class: string;

	let modal: HTMLDialogElement;
	onMount(() => {
		modal = document.getElementById("del-variable") as HTMLDialogElement;
	});

	let confirmKey: string = "";

	let error: string | undefined = undefined;
	let loading: boolean = false;

	async function handleSubmit() {
		try {
			if (!repoID) throw new Error("Repository ID not found");
			loading = true;
			await handleDelete(repoID, variableID);
			modal?.close();
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}
</script>

<button
	class={_class}
	on:click={() => modal?.showModal()}
>
	<i class="fa-solid fa-trash fa-sm" />
</button>

<dialog
	id="del-variable"
	class="modal"
>
	<form
		method="dialog"
		class="modal-box"
	>
		<div>
			<h3 class="mb-7 text-lg font-bold">Delete repository</h3>
			<p class="mb-4">
				You're about to delete the environment variable
				<span class="font-bold">
					{variableKey}
				</span>
			</p>
			<div class="alert alert-error mb-5">
				<i class="fa-solid fa-triangle-exclamation" />
				<p class="">
					<span class="font-bold">Warning:</span> this action is not reversible.
				</p>
			</div>

			<p class="mb-2">
				Enter
				<span class="font-bold">`{variableKey}`</span>
				to delete the variable:
			</p>

			<div class="flex h-[70px] flex-col items-start justify-start">
				<input
					type="text"
					class="input input-bordered w-full"
					bind:value={confirmKey}
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
					Close
				</button>
				<button
					class="btn btn-error w-44"
					type="submit"
					disabled={confirmKey !== variableKey}
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
		id="del-var-close"
		method="dialog"
		class="modal-backdrop"
	>
		<button id="delete-var-close" />
	</form>
</dialog>
