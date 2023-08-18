<script lang="ts">
	import { onMount } from "svelte";
	import type { NewVariable } from "./store";
	import classNames from "classnames";
	import { writeNewVariable } from "./requests";

	export let repoID: number;
	export let writeAccess: boolean;

	type DialogElement = HTMLDialogElement | null;
	let modal: DialogElement | null;
	onMount(() => {
		modal = document.querySelector("#new-var");
	});

	const v: NewVariable = {
		key: "",
		value: "",
	};

	let loading: boolean = false;
	let error: string | undefined;
	async function handleSubmit(e: Event) {
		e.preventDefault();
		try {
			loading = true;
			await writeNewVariable(repoID, v);
			modal?.close();
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}
</script>

<button
	class="btn btn-primary text-xs"
	disabled={!writeAccess}
	on:click={() => modal?.showModal()}
>
	<i class="fa-solid fa-plus" />
	Add
</button>

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
			<button
				class="btn btn-ghost w-28"
				type="button"
				on:click={() => modal?.close()}
			>
				Cancel
			</button>
			<button
				class={classNames(["btn btn-primary w-28"])}
				type="submit"
				disabled={loading}
			>
				{#if loading}
					<span class="loading loading-sm" />
				{:else}
					Submit
				{/if}
			</button>
		</div>
	</form>
	<form id="new-var-close" method="dialog" class="modal-backdrop">
		<button id="new-var-close" />
	</form>
</dialog>
