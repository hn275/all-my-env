<script lang="ts">
	import { afterUpdate, onMount } from "svelte";
	import type { NewVariable, RepositoryEnv } from "../store";
	import cn from "classnames";
	import { writeNewVariable } from "../requests";
	import { store } from "../store";

	let repo: RepositoryEnv;
	$: repo = $store;

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
		if (!repo.repoID) throw new Error("repository not found.");
		e.preventDefault();
		try {
			loading = true;
			await writeNewVariable(repo.repoID, v);
			modal?.close();
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	function handleOpen() {
		v.key = "";
		v.value = "";
		modal?.showModal();
	}

	afterUpdate(() => {
		if (!error) return;
		setTimeout(() => {
			error = undefined;
		}, 3000);
	});
</script>

<button
	class="btn btn-primary text-xs"
	disabled={!$store.write_access}
	on:click={handleOpen}
>
	<i class="fa-solid fa-plus" />
	Add
</button>

<dialog
	id="new-var"
	class="modal"
>
	<form
		method="dialog"
		class="modal-box"
		on:submit={handleSubmit}
	>
		<h3 class="text-gradient mb-3 text-lg font-bold">New variable</h3>
		<div>
			<div class="form-control relative mb-5">
				<label
					for="key"
					class="label"
				>
					Key*
				</label>
				<input
					required
					id="key"
					type="text"
					placeholder="foo"
					class="input input-bordered w-full"
					bind:value={v.key}
				/>
				{#if error}
					<p class="text-error absolute -bottom-5 left-1 text-xs">
						{error}
					</p>
				{/if}
			</div>

			<div>
				<label
					class="label"
					for="value"
				>
					Value*
				</label>
				<input
					required
					id="value"
					type="text"
					placeholder="bar"
					class="input input-bordered w-full"
					bind:value={v.value}
				/>
			</div>
		</div>

		<div class="mt-8 flex items-center justify-center gap-4">
			<button
				class="btn btn-ghost w-28"
				type="button"
				on:click={() => modal?.close()}
			>
				Cancel
			</button>
			<button
				class={cn(["btn btn-primary w-28"])}
				type="submit"
				disabled={loading || v.key === "" || v.value === ""}
			>
				{#if loading}
					<span class="loading loading-sm" />
				{:else}
					Submit
				{/if}
			</button>
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
