<script lang="ts">
	import Spinner from "@components/spinner.svelte";
	import { AuthStore, type User } from "@lib/auth";
	import { onMount } from "svelte";
	import cx from "classnames";
	import {
		PUBLIC_GITHUB_CLIENT_ID,
		PUBLIC_NODE_ENV,
	} from "$env/static/public";
	import Github from "@assets/github.svelte";

	export let loading: boolean;
	let oauth: string = "/dash";
	let user: User | undefined;
	onMount(() => {
		user = AuthStore.user();
		if (user) return;
		const isProd: boolean = PUBLIC_NODE_ENV === "production";
		const http = isProd ? "https://" : "http://";
		const redirect_uri = http + window.location.host + "/auth";
		const client_id = PUBLIC_GITHUB_CLIENT_ID;
		const scope = "repo user read:org";
		const p = new URLSearchParams({ client_id, redirect_uri, scope });
		oauth = "https://github.com/login/oauth/authorize?" + p.toString();
	});
	$: loggedIn = oauth === "/dash";
</script>

<div class="w-[14ch] overflow-x-hidden">
	{#if !loggedIn}
		<a
			href={oauth}
			class={cx([
				"btn normal-case text-dark-100 w-full h-full",
				loggedIn ? "btn-outline border-main text-main" : "btn-primary",
				{ "btn-disabled cursor-default": loading },
			])}
		>
			{#if loading}
				<Spinner />
			{:else}
				Sign in
				<span>
					<Github />
				</span>
			{/if}
		</a>
	{:else}
		<p class="font-semibold text-main text-gradient">
			Welcome back, {user?.name}
		</p>
	{/if}
</div>
