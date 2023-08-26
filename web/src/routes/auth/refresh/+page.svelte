<script lang="ts">
	import { AuthStore } from "@lib/auth";
	import type { User } from "@lib/auth";
	import { onMount } from "svelte";
	import { refreshSession } from "./requests";

	let url: string;
	let r: string;
	let data: EnvHub.Response<User>;
	onMount(async () => {
		r = AuthStore.redirectUrl() ?? "not found";
		url = window.location.href;
		try {
			data = await refreshSession();
			const red = AuthStore.redirectUrl();
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
