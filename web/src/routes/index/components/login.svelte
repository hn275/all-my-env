<script lang="ts">
	import Spinner from "@components/spinner.svelte";
	import { AuthStore, type User } from "@lib/auth";
	import { onMount } from "svelte";
	import {
		PUBLIC_GITHUB_CLIENT_ID,
		PUBLIC_NODE_ENV,
	} from "$env/static/public";
	import Github from "../assets/github.svelte";
	import Logout from "../assets/logout.svg";
	import Files from "../assets/files.svg";
	import { logout } from "./requests";
	import cx from "classnames";

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
	let loggedIn: boolean;

	let dropdownOpen: boolean = false;
	const handleDropdown = () => (dropdownOpen = !dropdownOpen);

	let logoutLoading: boolean = false;
	async function handleLogout(): Promise<void> {
		try {
			logoutLoading = true;
			await logout();
			AuthStore.logout();
			user = undefined;
		} catch (e) {
			console.error(e);
		} finally {
			logoutLoading = false;
		}
	}

	$: loggedIn = user !== undefined;
</script>

<div class="w-[14ch] overflow-x-hidden">
	{#if !loggedIn}
		<a
			href={oauth}
			class={cx([
				"btn text-dark-100 h-full w-full normal-case",
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
		<div class="relative">
			<button class="w-full hover:text-light transition" on:click={handleDropdown}>
				<img
					src={user?.avatar_url}
					alt={user?.name}
					class="rounded-full w-7 mb-2 inline mr-2"
				/>
				<span class="font-normal">
					{user?.name}
				</span>
			</button>

			<ul
				class={cx([
					"user-dropdown transition-all",
					"bg-dark-100 fixed right-3 top-full flex h-0 w-max -translate-y-[7px]",
					"flex-col items-start justify-center overflow-y-clip rounded-md",
					{ "h-[90px]": dropdownOpen },
				])}
			>
				<li>
					<a href="/dash">
						<img src={Files} alt="files" role="presentation" />
						<span>.env</span>
					</a>
				</li>

				<li>
					<button on:click={handleLogout}>
						{#if logoutLoading}
							loading
						{:else}
							<img src={Logout} alt="files" role="presentation" />
						{/if}
						<span>Log out</span>
					</button>
				</li>
			</ul>
		</div>
	{/if}
</div>

<style lang="postcss">
	.user-dropdown li {
		@apply w-max;
	}

	.user-dropdown button,
	.user-dropdown a {
		@apply flex justify-start items-center w-[125px] px-4 h-[40px];
	}

	.user-dropdown img {
		@apply w-4 mr-4;
	}
	.user-dropdown span {
		@apply font-normal;
	}

	.user-dropdown button:hover span,
	.user-dropdown a:hover span {
		@apply underline;
	}
</style>
