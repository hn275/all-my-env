<script lang="ts">
	import { onMount } from "svelte";
	import Nav from "./nav.svelte";
	import { auth, refresh } from "./auth";
	import type { User } from "$lib";
	import { TokenStorage, UserStorage } from "$lib/storage";

	let loading: boolean = false;
	let loadingMsg: string | undefined = undefined;
	let err: string | undefined;
	onMount(async () => {
		const url = new URL(window.location.href);
		const code: string | null = url.searchParams.get("code");
		let user: User | undefined = UserStorage.get();
		try {
			loading = true;
			if (!code) {
				if (TokenStorage.refreshed() || !user) return;
				loadingMsg = "refreshing token...";
				user = await refresh(user.access_token);
			} else {
				loadingMsg = "signing in...";
				user = await auth(code);
			}
		if (!user) return;
			UserStorage.set(user);
			TokenStorage.setRefreshed();
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
			loadingMsg = undefined;
		}
	});
</script>

<Nav {loadingMsg} />
<h1 class="text-light">Hello world</h1>

{#if loading}
	<p>Loading</p>
{/if}
