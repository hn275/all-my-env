<script lang="ts">
	import { onMount } from "svelte";
	import Nav from "./index/components/nav.svelte";
	import { Auth } from "@lib/auth";
	import type { User } from "@lib/auth";
	import Faq from "./index/components/faq.svelte";
	import Login from "./index/components/login.svelte";

	let user: User | undefined;
	onMount(() => {
		user = Auth.user();
	});
</script>

<Nav />
<main class="relative mx-auto max-w-7xl scroll-m-28 overflow-x-hidden">
	<section
		class="mb-24 mt-64 flex flex-col items-center justify-center gap-6"
	>
		<h1
			class="text-gradient font-accent text-center text-3xl font-bold uppercase md:text-4xl"
		>
			EFFORTLESS SECRETS MANAGEMENT
		</h1>

		<p class="w-[25ch] text-center font-semibold md:w-max md:text-lg">
			The solution for managing your
			<span class="hero-text-underline">variables</span>.
		</p>

		<div
			class="flex flex-col items-center justify-center gap-3 md:flex-row md:gap-6"
		>
			{#if user}
				<a
					href="/dashboard"
					class="btn btn-primary"
				>
					Get started
				</a>
			{:else}
				<Login />
			{/if}
		</div>
		<div>
			<p class="badge badge-outline badge-info text-xs">
				EnvHub is in beta.
			</p>
		</div>
	</section>

	<!-- FEATURES -->
	<section class="my-24 h-max">
		<ul
			class="mx-auto grid max-w-3xl grid-rows-3 gap-8 px-5 md:grid-cols-3 md:grid-rows-1"
		>
			<li class="feature-card">
				<div class="icon">
					<i class="fa-brands fa-github fa-2xl" />
				</div>
				<div>
					<h3 class="text-gradient">Easy Integration</h3>
					<p>
						EnvHub works closely with GitHub API, meaning changes
						are updated in real time.
					</p>
				</div>
			</li>

			<li class="feature-card">
				<div class="icon">
					<i class="fa-solid fa-key fa-2xl"></i>
				</div>
				<div>
					<h3 class="text-gradient">Secure</h3>
					<p>
						Your secrets are safe with us! Connections are
						encrypted, variables are ciphered.
					</p>
				</div>
			</li>

			<li class="feature-card">
				<div class="icon">
					<i class="fa-solid fa-gears fa-2xl"></i>
				</div>
				<div>
					<h3 class="text-gradient">Performance</h3>
					<p>
						EnvHub is built with Svelte and Golang, both offer
						wonderful DX, as well as performance stats.
					</p>
				</div>
			</li>
		</ul>
	</section>

	<!-- FAQ -->
	<section
		id="faqs"
		class="my-24 grid place-items-center"
	>
		<div class="mx-auto max-w-3xl">
			<h2 class="text-gradient mb-5 text-center text-2xl font-bold">
				Frequently Asked Questions
			</h2>

			<Faq />
		</div>
	</section>

	<!-- Support -->
	<section
		id="support"
		class="my-24"
	>
		<h2 class="text-gradient mb-3 text-center text-2xl font-bold">
			Support the maintainers
		</h2>

		<div
			class="mx-auto flex w-full max-w-[65ch] flex-col items-center gap-3"
		>
			<p class="text-center">
				EnvHub is 100% free! But there is still work in maintainence.
				<br />
				We appreciate what we can get.
			</p>

			<a
				href="/"
				class="btn btn-outline btn-accent w-max">Buy us a snack!</a
			>
		</div>
	</section>

	<section class="my-24 flex flex-col items-center gap-2">
		<p class="text-center font-bold">
			EnvHub is an
			<span class="hero-text-underline">open source</span>
			project.
		</p>

		<a
			href="https://github.com/hn275/envhub"
			target="_blank"
			class="link link-primary mx-auto font-normal"
		>
			Star us on GitHub!
		</a>
	</section>
</main>

<footer class="text-base-content/50 flex justify-center gap-5 py-5 text-sm">
	<a
		href="/terms"
		class="link hover:link-primary transition-all">Terms & Conditions</a
	>
	<a
		href="/privacy"
		class="link hover:link-primary transition-all">Privacy Agreement</a
	>
</footer>

<style lang="postcss">
	* {
		scroll-margin: 80px;
	}
	.hero-text-gradient {
		@apply bg-gradient-to-tr from-light to-main;
		background-clip: text;
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	span.hero-text-underline {
		position: relative;
		font-weight: 700;
		isolation: isolate;
	}

	span.hero-text-underline::after {
		content: "";
		position: absolute;
		bottom: 0;
		left: 0;
		width: 100%;
		height: 6px;
		border-radius: 4px;
		z-index: -1;
		@apply bg-accent/70;
	}

	.feature-card {
		@apply bg-neutral p-5 rounded-lg shadow-lg h-full w-full;
	}
	.feature-card h3 {
		@apply font-bold text-center mb-2 mt-3;
	}
	.feature-card .icon {
		@apply text-primary h-14 text-center grid place-items-center;
	}

	.feature-card p {
		@apply text-sm text-neutral-content text-start;
	}
</style>
