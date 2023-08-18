<script lang="ts">
	import { afterUpdate } from "svelte";

	export let created_at: string;
	export let updated_at: string;
	export let key: string;
	export let value: string;
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
		navigator.clipboard.writeText(`${key}="${value}"`).then(() => {
			copyOK = true;
		});
	}
	afterUpdate(() => {
		if (!copyOK) return;
		setTimeout(() => {
			copyOK = false;
		}, 2000);
	});
</script>

<tr>
	<td>{i + 1}</td>
	<td class="w-5">
		<p class="w-20 overflow-x-clip">
			{key}
		</p>
	</td>
	<td class="relative">
		<button
			on:click={copy}
			class="mr-2 text-light/30 hover:text-main/80 transition"
		>
			<i class="fa-regular fa-copy" />
		</button>
		<p class="inline">
			{value}
		</p>
	</td>
	<td>{createdAt}</td>
	<td>{updatedAt}</td>
</tr>

{#if copyOK}
	<div class="toast toast-start">
		<div class="alert alert-success flex justify-center">
			<p class="font-semibold text-center">Copied to clipboard!</p>
		</div>
	</div>
{/if}
