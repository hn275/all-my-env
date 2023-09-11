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

			<section class="join join-vertical px-5">
				<Faq header="What is the purpose of EnvHub?">
					<p>
						EnvHub offers an innovative solution for the efficient
						management of your critical environment variables. Our
						platform seamlessly integrates with GitHub, simplifying
						the task of sharing variables among collaborators.
					</p>
				</Faq>

				<Faq header="How does it work?">
					<p>
						EnvHub functionalities can be summarized in <span
							class="font-bold">3 steps</span
						>:
					</p>
					<ol class="ml-8 mt-2 list-decimal">
						<li>Authentication with GitHub.</li>
						<li>Create projects and their associated variables.</li>
						<li>Share it with your collaborators!</li>
					</ol>
				</Faq>

				<Faq header="Are my secrets safe with EnvHub?">
					<p>
						We exercise meticulouse care when handling your
						environment variables. Our server is hosted with <a
							href="https://www.cloudflare.com/learning/ssl/what-is-ssl/"
							target="_blank"
							class="link link-info font-semibold"
							>SSL Certificates</a
						>
						and your variables undergo encryption and decryption using
						the robust
						<a
							href="https://ieeexplore.ieee.org/document/7927078"
							target="_blank"
							class="link link-info font-semibold"
							>ChaCha20-Poly1305</a
						>. This choice reflects our commitment to maintaining a
						high-security standard while optimizing performance.
					</p>
				</Faq>

				<Faq header="Who built this app?">
					<p>
						EnvHub was initiated by computer science students from
						the
						<a
							href="https://uvic.ca"
							target="_blank"
							class="link link-info font-semibold"
							>University of Victoria</a
						>. Led by
						<a
							href="https://github.com/hn275"
							target="_blank"
							class="link link-info font-semibold">Hal</a
						>, a proactive student balancing academics with
						professional engagements, the project was born from the
						need for simplified variable sharing within teams, both
						in professional and collaborative settings.
					</p>
				</Faq>

				<Faq header="Why do I need to authenticate with push access?">
					<p>
						EnvHub <span class="font-bold">does not</span> initiate
						modifications to your GitHub account or repositories.
						Nevertheless, certain APIs we utilize necessitate push
						access, primarily for the purpose of verifying
						<a
							class="link link-info font-semibold"
							target="_blank"
							href="https://docs.github.com/en/rest/collaborators/collaborators?apiVersion=2022-11-28#list-repository-collaborators"
							>collaborators' status</a
						>.
					</p>
				</Faq>

				<Faq header="What's the tech stack?">
					<p>
						We prioritize performance excellence in our approach.
						For our front-end development, we employ
						<a
							href="https://kit.svelte.dev/"
							target="_blank"
							class="link link-info font-semibold">SvelteKit</a
						>
						combined with
						<span class="font-semibold">TypeScript</span>,
						delivering not only high efficiency but also robust type
						safety. Our UI components and color themes are
						meticulously crafted with
						<a
							href="https://daisyui.com/"
							class="link link-info font-semibold"
							target="_blank">DaisyUI</a
						>

						, ensuring a visually appealing and cohesive user
						experience.
					</p>
					<p class="mt-4">
						Behind the scenes, our application is powered by a
						<span class="font-semibold">Golang</span> server, a testament
						to our commitment to performance and reliability. Additionally,
						we place a strong emphasis on providing an exceptional development
						experience for our maintainers, recognizing the value of
						a smooth and efficient workflow.
					</p>
				</Faq>

				<Faq header="Why &quot;EnvHub&quot;?">
					<p>
						This question is best answered by the talented <a
							href="https://github.com/jaspreetks"
							target="_blank"
							class="link link-info font-semibold">Jaspreet</a
						> (who also designed this UI!)
					</p>
				</Faq>
			</section>
		</div>
	</section>

	<!-- Support -->
	<!--
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
-->

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

<!--
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
-->

<style lang="postcss">
	* {
		scroll-margin: 80px;
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
