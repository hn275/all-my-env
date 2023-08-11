<script lang="ts">
	import { onMount } from "svelte";
	import Nav from "./nav.svelte";
	import { signIn, refresh } from "./auth";
	import type { AuthStore, User } from "@lib/auth";
	import { authStore } from "@lib/auth";
	import Main from "@components/main.svelte";
	import { goto } from "$app/navigation";

	let loading: boolean = false;
	let err: string | undefined;
	onMount(async () => {
		const url = new URL(window.location.href);
		const code: string | null = url.searchParams.get("code");
		let user: User;
		const authState: AuthStore = authStore.get();

		try {
			loading = true;
			if (!code) {
				if (authState.tokenRefreshed || !authState.user) return;
				user = await refresh(authState.user.access_token);
			} else {
				user = await signIn(code);
				url.searchParams.delete("code");
				goto(url);
			}
			authStore.set({ tokenRefreshed: true, user });
		} catch (e) {
			err = (e as Error).message;
		} finally {
			loading = false;
		}
	});
</script>

<Nav {loading} />
<Main>
	<h1 class="text-light">Hello world</h1>

	{#if loading}
		<p>Loading</p>
	{/if}
</Main>
