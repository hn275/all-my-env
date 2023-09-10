<script lang="ts">
	import { Auth } from "@lib/auth";
	import { onMount } from "svelte";
	import { refreshSession } from "./requests";

	let url: string;
	let r: string;
	onMount(async () => {
		url = window.location.href;
		try {
			// data = await refreshSession();
			const red = Auth.redirectUrl();
			window.location.replace(red);
		} catch (e) {
			console.error(e);
		}
	});
</script>

<main>
	<p>
		refresh from: {url}
	</p>
	<p>to: {r}</p>
</main>
