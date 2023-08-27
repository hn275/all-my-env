<script lang="ts">
	import { createEventDispatcher, onMount } from "svelte";

	export let key: string;
	export let newKey: string;
	export let value: string;
	export let newValue: string;

	let modal: HTMLDialogElement;
	onMount(() => {
		modal = document.getElementById("cancel-modal") as HTMLDialogElement;
	});

	const dispatch = createEventDispatcher<{ undo: void }>();

	function handleOpen() {
		const sameKey: boolean = key === newKey;
		const sameValue: boolean = value === newValue;
		if (sameKey && sameValue) {
			dispatch("undo");
			return;
		}
		modal.showModal();
	}

	function handleReset() {
		dispatch("undo");
	}
</script>

<button
	on:click={handleOpen}
	class="btn-variable-utilities"
>
	<i class="fa-solid fa-xmark fa-sm" />
</button>

<dialog
	id="cancel-modal"
	class="modal"
>
	<form
		method="dialog"
		class="modal-box bg-dark-200"
	>
		<h3 class="mb-6 font-bold">Undo changes?</h3>

		<p class="mb-3">The followling changes will be discarded.</p>

		<div class="mb-3 flex items-center gap-3">
			<h6 class="w-[4ch] font-semibold">Key</h6>
			<p class="bg-dark-100 flex items-center gap-3 rounded-md p-3">
				{#if key !== newKey}
					<span>{key}</span>
					<i class="fa-solid fa-arrow-right" />
					<span>{newKey}</span>
				{:else}
					<span>{key}</span>
				{/if}
			</p>
		</div>

		<div class="flex items-center gap-3">
			<h6 class="w-[4ch] font-semibold">Value</h6>
			<p class="bg-dark-100 flex items-center gap-3 rounded-md p-3">
				{#if value !== newValue}
					<span>{value}</span>
					<i class="fa-solid fa-arrow-right" />
					<span>{newValue}</span>
				{:else}
					<span>{value}</span>
				{/if}
			</p>
		</div>

		<div class="mt-6 flex items-center justify-end gap-3">
			<button
				class="btn btn-ghost"
				type="button"
				on:click={() => modal.close()}
			>
				Cancel
			</button>
			<button
				class="btn btn-error"
				type="submit"
				on:click|preventDefault={handleReset}
			>
				Confirm
			</button>
		</div>
	</form>
	<form
		id="edit-var-close"
		method="dialog"
		class="modal-backdrop"
	>
		<button id="cancel-btn" />
	</form>
</dialog>
