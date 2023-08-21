<script lang="ts">
	import { afterUpdate } from "svelte";
	import cx from "classnames";
	import Row from "./row.svelte";
	import { deleteVariable } from "../services";
	import type { Variable } from "../store";

	export let created_at: string;
	export let updated_at: string;
	export let key: string;
	export let value: string;
	export let id: string;
	export let i: number;

	function formatTime(d: Date): string {
		let dt = d.toLocaleDateString() + " ";
		dt += d.getHours().toString() + ":";
		dt += d.getMinutes().toString().padStart(2, "0");
		return dt;
	}

	const createdAt = formatTime(new Date(created_at));
	const updatedAt = formatTime(new Date(updated_at));

	let copyOK: boolean = false;
	function copy() {
		navigator.clipboard
			.writeText(`${key}="${value}"`)
			.then(() => (copyOK = true));
	}

	afterUpdate(() => {
		if (!copyOK) return;
		setTimeout(() => {
			copyOK = false;
		}, 2000);
	});

	let varKey: string = key;
	let varValue: string = value;
	let editMode: boolean = false;

	let confirmCancel: boolean = false;
	function reset() {
		editMode = false;
		confirmCancel = false;
		varKey = key;
		varValue = value;
	}

	let saveAble: boolean;
	$: saveAble = !(varKey === key) || !(varValue === value);
	let editLoading: boolean = false;
	async function handleSubmit() {}

	function handleDelete() {
		const v: Variable = { id, key, value, updated_at, created_at };
		deleteVariable(v);
	}
</script>

<Row className="group">
	<div class="group flex items-center justify-start gap-3">
		<div
			class={cx([
				"transitio flex flex-row gap-1 transition-all",
				{ "-ml-12 opacity-0": !editMode },
				"group-hover:ml-0 group-hover:opacity-100",
			])}
		>
			{#if !editMode}
				<!-- delete button -->
				<button class="delete" on:click={handleDelete}>
					<i class="fa-solid fa-trash fa-sm" />
				</button>
				<!-- edit button -->
				<button on:click={() => (editMode = true)}>
					<i class="fa-regular fa-pen-to-square fa-sm" />
				</button>
				<!-- copy button -->
				<button on:click={copy} class="button">
					<i class="fa-regular fa-copy fa-sm" />
				</button>
			{:else}
				<!-- cancel button -->
				<button on:click={reset}>
					<i class="fa-solid fa-xmark fa-sm" />
				</button>
				<!-- save button -->
				<button
					on:click={handleSubmit}
					disabled={editLoading || !saveAble}
					class="save"
				>
					<i class="fa-solid fa-check fa-sm" />
				</button>
			{/if}
		</div>
		{#if !editMode}
			<p
				class="text-light/40 relative bottom-[2px] flex-grow self-end text-sm group-hover:opacity-0"
			>
				{i + 1}.
			</p>
		{/if}
	</div>

	<input class="bg-transparent transition-all" bind:value={varKey} disabled={!editMode} />

	<div class="relative">
		<input
			class="text-main w-full bg-transparent font-semibold transition-all"
			bind:value={varValue}
			disabled={!editMode}
		/>
	</div>

	<p class="text-light/70 text-sm">{createdAt}</p>

	<p class="text-light/70 text-sm">{updatedAt}</p>

	{#if copyOK}
		<div class="toast toast-start">
			<div class="alert alert-success flex justify-center">
				<p class="text-center font-normal">Copied to clipboard!</p>
			</div>
		</div>
	{/if}

	{#if confirmCancel}
		<div>confirm cancle</div>
	{/if}
</Row>

<style lang="postcss">
	input,
	p {
		@apply text-ellipsis;
		@apply p-2;
	}

	input:disabled {
		@apply bg-transparent;
	}

	input {
		@apply bg-dark-100 rounded-md;
	}

	input:focus {
		@apply border border-light/60;
	}

	button {
		@apply w-6 h-6 rounded-md hover:bg-light/10;
		@apply text-light/50 hover:text-light transition;
	}

	button.save {
		@apply text-light bg-main;
		@apply hover:brightness-110 hover:bg-main;
	}

	button.delete {
		@apply hover:bg-error hover:text-dark-200;
	}

	button:disabled {
		@apply opacity-20;
	}
</style>
