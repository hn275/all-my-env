<script lang="ts">
	import { onMount } from "svelte";

	export let id: string;
	export let heading: string;

	function handleClose() {
		document.querySelector("body")?.classList.remove("overflow-y-hidden");
	}

	function eventSub(e: Event) {
		const esc = (e as KeyboardEvent).key === "Escape";
		if (!esc) return;
		document.querySelector("body")?.classList.remove("overflow-y-hidden");
	}

	onMount(() => {
		document.addEventListener("keydown", eventSub);
		return () => document.removeEventListener("keydown", eventSub);
	});
</script>

<dialog
	{id}
	class="modal"
>
	<form
		method="dialog"
		class="modal-box"
	>
		<h2 class="text-gradient mb-6 text-lg font-bold">{heading}</h2>
		<slot />
	</form>
	<form
		method="dialog"
		class="modal-backdrop"
	>
		<button on:click={handleClose}></button>
	</form>
</dialog>
