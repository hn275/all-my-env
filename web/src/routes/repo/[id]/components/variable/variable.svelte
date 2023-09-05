<script lang="ts">
	import { afterUpdate } from "svelte";
	import cx from "classnames";
	import Row from "./row.svelte";
	import DeleteVariable from "./delete-variable.svelte";
	import ConfirmEdit from "./confirm-edit.svelte";
	import ConfirmCancel from "./confirm-cancel.svelte";
	import type { RepositoryEnv } from "../../store";
	import { store } from "../../store";
	import { formatTime } from "../../services";

	export let created_at: string;
	export let updated_at: string;
	export let key: string;
	export let value: string;
	export let id: string;
	export let i: number;

	let state: RepositoryEnv;
	$: state = $store;

	let createdAt: string = formatTime(new Date(created_at));

	let updatedAt: string;
	$: updatedAt = formatTime(new Date(updated_at));

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

	let newKey: string = key;
	let newValue: string = value;
	let editMode: boolean = false;

	function handleReset() {
		editMode = false;
		newKey = key;
		newValue = value;
	}

	let saveAble: boolean;
	$: saveAble = !(newKey === key) || !(newValue === value);

	function handleEditOK() {
		editMode = false;
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
				<DeleteVariable
					repoID={state.repoID}
					variableID={id}
					variableKey={key}
				/>
				<!-- edit button -->
				<button
					on:click={() => (editMode = true)}
					class="btn-variable-utilities"
				>
					<i class="fa-regular fa-pen-to-square fa-sm" />
				</button>
				<!-- copy button -->
				<button
					on:click={copy}
					class="btn-variable-utilities"
				>
					<i class="fa-regular fa-copy fa-sm" />
				</button>
			{:else}
				<!-- cancel button -->
				<ConfirmCancel
					{key}
					{newKey}
					{value}
					{newValue}
					on:undo={handleReset}
				/>
				<!-- save button -->
				<ConfirmEdit
					on:success={handleEditOK}
					{id}
					{saveAble}
					{key}
					{newKey}
					{value}
					{newValue}
				/>
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

	<input
		class="bg-transparent transition-all"
		bind:value={newKey}
		disabled={!editMode}
	/>

	<div class="relative">
		<input
			class="text-primary w-full bg-transparent font-semibold transition-all"
			bind:value={newValue}
			disabled={!editMode}
		/>
	</div>

	<p class="text-base-content/70 text-sm">{createdAt}</p>

	<p class="text-base-content/70 text-sm">{updatedAt}</p>

	{#if copyOK}
		<div class="toast toast-start">
			<div class="alert alert-success flex justify-center">
				<p class="text-center font-normal">Copied to clipboard!</p>
			</div>
		</div>
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

	button:disabled {
		@apply opacity-20;
	}
</style>
