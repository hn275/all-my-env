<script lang="ts">
	import { onMount } from "svelte";
	import Nav from "./index/components/nav.svelte";
	import { refresh } from "./auth";
	import { AuthStore } from "@lib/auth";
	import type { User } from "@lib/auth";
	import cx from "classnames";
	import LogInBtn from "./index/components/login.svelte";
	import StarUs from "./index/components/star.svelte";
	import Features from "./index/components/features.svelte";

	let loading: boolean = false;
	let err: string | undefined;
	let user: User | undefined;
	onMount(async () => {
		try {
			user = AuthStore.user();
			if (!user) return;
			const refreshed = AuthStore.sessionRefreshed();
			if (refreshed) return;
			loading = true;
			user = await refresh();
			AuthStore.login(user);
			AuthStore.refreshed();
		} catch (e) {
			err = (e as Error).message;
			AuthStore.logout();
		} finally {
			loading = false;
		}
	});
</script>

<Nav {loading} />
<main
	class="relative mx-auto h-max min-h-[calc(100vh-65px)] max-w-7xl overflow-x-hidden"
>
	<section
		class="flex h-[calc(100vh-65px)] flex-col items-center justify-center gap-6"
	>
		<h1
			class="hero-text-gradient font-accent text-center text-3xl font-bold uppercase md:text-4xl"
		>
			effortless secrets management
		</h1>

		<p class="w-[25ch] text-center font-semibold md:w-max md:text-lg">
			The ultimate{" "}
			<span class="hero-text-underline">open-source</span>
			&nbsp; solution for managing your&nbsp;
			<span class="hero-text-underline">variables</span>.
		</p>

		<div
			class="flex flex-col items-center justify-center gap-3 md:flex-row md:gap-6"
		>
			<a
				href="#get-started"
				class="btn btn-outline border-main text-main w-[14ch]"
			>
				Get started
			</a>

			{#if user}
				<a
					href="/dashboard"
					class="btn btn-primary w-[14ch]"
				>
					Dashboard
				</a>
			{:else}
				<LogInBtn {loading} />
			{/if}
		</div>
	</section>

	<section
		class="mx-auto mb-40 w-full max-w-4xl"
		id="get-started"
	>
		<h2
			class="font-accent hero-text-gradient text-center text-4xl font-bold md:ml-12 md:text-left"
		>
			How it works?
		</h2>
		<Features />
	</section>
</main>
<footer
	class={cx([
		"relative mx-auto grid max-w-screen-2xl grid-cols-3 place-items-center",
		"text-light/40 py-5 text-sm",
	])}
>
	<div>
		<p>EnvHub</p>
	</div>

	<StarUs />

	<div class="flex gap-5">
		<a
			href="/terms"
			class="hover:text-light">Terms & Conditions</a
		>
		<a
			href="/privacy"
			class="hover:text-light">Privacy Agreement</a
		>
	</div>
</footer>
<div class="hero-graphic-blur main" />
<div class="hero-graphic-blur accent" />

<style lang="postcss">
	.hero-text-gradient {
		@apply bg-gradient-to-tr from-light to-main;
		background-clip: text;
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	span.hero-text-underline {
		position: relative;
		font-weight: 700;
	}

	span.hero-text-underline::after {
		content: "";
		position: absolute;
		bottom: 0;
		left: 0;
		width: 100%;
		height: 6px;
		border-radius: 4px;
		background: red;
		z-index: -1;
		@apply bg-accent-fuchsia-100/80;
	}

	.hero-graphic-blur {
		position: absolute;
		bottom: 0;
		right: 0;
		filter: blur(150px);
		transform-origin: center;
	}

	.hero-graphic-blur.main {
		height: 325px;
		width: 500px;
		@apply bg-main/30;
		z-index: -20;
		transform: scale(0.5) translate(250px, 200px);
	}

	.hero-graphic-blur.accent {
		height: 400px;
		width: 470px;
		@apply bg-accent-fuchsia-100/30;
		z-index: -30;
		transform: scale(0.5) translate(calc(40% + 200px), -30px);
	}

	@media only screen and (min-width: 768px) {
		.hero-graphic-blur.main {
			transform: translate(0, 0px);
		}

		.hero-graphic-blur.accent {
			transform: translate(40%, -150px);
		}
	}
</style>
