<script lang="ts">
	import { urlParam, urlRedirect } from "@lib/url";
	import { signIn } from "./auth";
	import { onMount } from "svelte";
	import { AuthStore } from "@lib/auth";
	import Spinner from "@components/spinner.svelte";

	let error: string | null;
	let code: string | null;
	let loading: boolean = true;
	onMount(async () => {
		error = urlParam("error");
		code = urlParam("code");
		try {
			if (!code) return;
			const user = await signIn(code);
			AuthStore.login(user);
			urlRedirect("/dashboard");
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	});
</script>

<main class="flex min-h-[100vh] items-center justify-center">
	{#if error}
		<section class="flex flex-col items-center gap-5 p-9">
			<h2 class="text-accent-fuchsia-100 text-3xl font-semibold">
				Oops, something went wrong &#128532;
			</h2>
			<p>
				{error}
			</p>
			<a href="/" class="btn btn-outline">Back to home page</a>
		</section>
	{:else if loading}
		<Spinner class="stroke-dark-200" />
		<p>Authenticating...</p>
	{/if}
</main>
