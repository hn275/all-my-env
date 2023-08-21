<script lang="ts">
	import { store } from "../store";
	import { cancelDelete, handleDelete } from "../services";
	import { onMount } from "svelte";

	$: state = $store;
	let confirmKey: string = "";

	onMount(() => {
		confirmKey = "";
		function sub(e: KeyboardEvent) {
			switch (e.key) {
				case "Escape":
					cancelDelete();
					return;
				default:
					break;
			}
		}
		document.addEventListener("keydown", sub);
		return () => document.removeEventListener("keydown", sub);
	});

	$: disabled = confirmKey !== state.deleteVariable?.key || loading;
	let loading: boolean = false;
	let error: string | undefined;
	async function onSubmit() {
		if (!state.deleteVariable || !state.repoID || disabled) return;
		try {
			loading = true;
			await handleDelete(state.repoID, state.deleteVariable.id);
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}
</script>

{#if state.deleteVariable}
	<form
		id="#delete-modal"
		class="bg-dark-100/75 fixed bottom-0 left-0 right-0 top-0 grid place-items-center backdrop-blur transition-all"
		on:submit|preventDefault={onSubmit}
	>
		<div class="bg-dark-200 modal-box">
			<h3 class="mb-7 text-lg font-bold">Delete a variable</h3>
			<p class="mb-4">
				You're about to delete the variable
				<span class="font-bold">
					`{state.deleteVariable.key}`
				</span>.
			</p>

			<div class="alert alert-error mb-5">
				<i class="fa-solid fa-triangle-exclamation" />
				<p class="">
					<span class="font-bold">Warning:</span> this action is not reversible.
				</p>
			</div>

			<p class="mb-2">
				Enter the variable key
				<span class="font-bold">`{state.deleteVariable.key}`</span> to continue:
			</p>

			<div class="flex h-[70px] flex-col items-start justify-start">
				<input
					type="text"
					bind:value={confirmKey}
					class="input input-bordered w-full"
					placeholder={state.deleteVariable.key}
				/>
				{#if error}
					<p class="text-error text-xs">{error}</p>
				{/if}
			</div>
			<div class="modal-action">
				<button
					class="btn btn-ghost"
					on:click={cancelDelete}
					type="button"
				>
					Close
				</button>
				<button class="btn btn-primary w-44" {disabled} type="submit">
					{#if loading}
						<span class="loading" />
					{:else}
						Delete my variable
					{/if}
				</button>
			</div>
		</div>
	</form>
{/if}
