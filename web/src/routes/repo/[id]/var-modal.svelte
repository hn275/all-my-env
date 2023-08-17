<script lang="ts">
	import { afterUpdate, createEventDispatcher } from "svelte";
	import type { NewVariable } from "./store";
	import classNames from "classnames";

	export let open: boolean;
	export let loading: boolean;

	const dispatch = createEventDispatcher<{ write: NewVariable; close: void }>();
	const v: NewVariable = {
		key: "",
		value: "",
	};

	function handleSubmit(e: Event) {
		e.preventDefault();
		dispatch("write", v);
	}

	function onClose() {
		dispatch("close");
	}
	type DialogElement = HTMLDialogElement | null;
	afterUpdate(() => {
		if (document) {
			const modal: DialogElement | null = document.querySelector("#new-var");
			open ? modal?.showModal() : modal?.close() && onClose();
		}
	});
</script>

<dialog id="new-var" class="modal">
	<form method="dialog" class="modal-box bg-dark-200" on:submit={handleSubmit}>
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
			<button class="btn btn-ghost w-28" type="button" on:click={onClose}>
				Cancel
			</button>
			<button class={classNames(["btn btn-primary w-28"])} type="submit">
				{#if loading}
					<span class="loading loading-sm" />
				{:else}
					Submit
				{/if}
			</button>
		</div>
	</form>
	<form id="new-var-close" method="dialog" class="modal-backdrop">
		<button id="new-var-close" on:click={onClose} />
	</form>
</dialog>
