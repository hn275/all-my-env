<script lang="ts">
	import { store, type Variable } from "../store";
	import { cancelDelete } from "../services";
	import { onDestroy, onMount } from "svelte";

	let variable: Variable | undefined;
	$: variable = $store.state.deleteVariable;
    let confirmKey: string = "";

    onMount(() => {
      document.addEventListener("keydown", (e) => {
        if (e.key !== "Escape") return
        cancelDelete();
        confirmKey = "";
      })
    })


    onDestroy(() => {
      confirmKey = "";
    })
</script>

{#if variable}
	<form
		id="#delete-modal"
		class="bg-dark-100/75 fixed bottom-0 left-0 right-0 top-0 grid place-items-center backdrop-blur transition-all"
	>
		<div class="bg-dark-200 modal-box">
			<h3 class="mb-7 text-lg font-bold">Delete a variable</h3>
			<p class="mb-4">
				You're about to delete the variable
				<span class="font-bold">
					`{variable.key}`
				</span>.
			</p>

			<div class="alert alert-error mb-5">
				<i class="fa-solid fa-triangle-exclamation" />
				<p class="">
					<span class="font-bold">Warning:</span> this action is not reversible.
				</p>
			</div>

			<p>
				Enter the variable key
				<span class="font-bold">`{variable.key}`</span> to continue:
			</p>

			<input
				type="text"
				bind:value={confirmKey}
				class="input input-bordered w-full"
				placeholder={variable.key}
			/>
			<div class="modal-action">
				<button class="btn btn-ghost" on:click={cancelDelete}
					>Close</button
				>
				<button
					class="btn btn-primary"
					disabled={confirmKey !== variable.key}
					>Delete my variable</button
				>
			</div>
		</div>
	</form>
{/if}
