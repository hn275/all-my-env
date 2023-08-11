<script lang="ts">
	import { onMount } from "svelte";
	import Nav from "./nav.svelte";
	import { refresh } from "./auth";
	import { AuthStore } from "@lib/auth";

	let loading: boolean = false;
	let err: string | undefined;
	onMount(async () => {
		try {
			let user = AuthStore.user();
			if (AuthStore.sessionRefreshed() || !user) return;
			loading = true;
			user = await refresh(user.access_token);
			AuthStore.login(user);
			AuthStore.refreshed();
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	});
</script>

<Nav {loading} />
<h1 class="text-light">Hello world</h1>

{#if loading}
	<p>Loading</p>
{/if}
