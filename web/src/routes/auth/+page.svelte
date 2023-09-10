<script lang="ts">
	import { urlParam } from "@lib/url";
	import { signIn } from "./auth";
	import { onMount } from "svelte";
	import { Auth } from "@lib/auth";

	let error: string | null;
	let code: string | null;
	let loading: boolean = true;
	onMount(async () => {
		error = decodeURIComponent(urlParam("error") ?? "");
		code = urlParam("code");
		try {
			if (!code) {
				return;
			}
			const user = await signIn(code);
			Auth.login(user);
			Auth.refresh();
			window.location.replace("/dashboard");
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
			<a
				href="/"
				class="btn btn-outline btn-primary">Back to home page</a
			>
		</section>
	{:else if loading}
		<div class="text-primary text-center">
			<span class="loading loading-lg loading-ring"></span>
			<p class="text-info text-lg font-semibold">
				Authenticating with GitHub
			</p>
		</div>
	{/if}
</main>
