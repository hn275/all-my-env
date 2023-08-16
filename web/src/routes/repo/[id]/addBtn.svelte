<script lang="ts">
	import { createEventDispatcher } from "svelte";
	import type { NewVariable } from "./store";

	export let writeAccess: boolean;
	export let loading: boolean;
	export let error: string | undefined;

	function openModal(): void {
		type DialogElement = HTMLDialogElement | null;
		const modal: DialogElement = document.querySelector("#new-var");
		modal?.showModal();
	}

	const dispatch = createEventDispatcher<{write: NewVariable}>();
	const v: NewVariable = {
		key: "",
		value: "",
	};
	function handleSubmit(e: Event): void {
		e.preventDefault();
		dispatch("write", v);
		// close modal
		type ButtonElement = HTMLButtonElement | null;
		const modal: ButtonElement = document.querySelector("#new-var-close");
		modal?.click();
	}
</script>

<button
	class="btn btn-primary text-xs"
	disabled={!writeAccess}
	on:click={openModal}
>
	<i class="fa-solid fa-plus" />
	Add
</button>

<dialog id="new-var" class="modal">
	<form method="dialog" class="modal-box bg-dark-200">
		<h3 class="text-main text-lg font-bold">New variable</h3>
		<div class="">
			<div>
				<label for="key" class="label">Key</label>
				<input
					id="key"
					type="text"
					placeholder="foo"
					class="input input-bordered w-full"
					bind:value={v.key}
				/>
			</div>

			<div>
				<label class="label" for="value">Value</label>
				<input
					id="value"
					type="text"
					placeholder="bar"
					class="input input-bordered w-full"
					bind:value={v.value}
				/>
			</div>
		</div>

		<div class="mt-5 flex items-center justify-center gap-4">
			<button class="btn btn-ghost">Cancel</button>
			<button class="btn btn-primary" on:click={handleSubmit}>Submit</button>
		</div>
	</form>
	<form id="new-var-close" method="dialog" class="modal-backdrop">
		<button id="new-var-close" />
	</form>
</dialog>
