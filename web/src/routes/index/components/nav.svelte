<script lang="ts">
	import cx from "classnames";
	import Logo from "@assets/logo.svg";
	import Hamburger from "../assets/hamburger.svelte";
	import Dashboard from "./dashboard.svelte";

	let show: boolean = false;
	const toggleOpen: () => void = () => (show = !show);

	type NavMap = {
		href: string;
		text: string;
	};
	const navMap: Array<NavMap> = [
		{ href: "/#faqs", text: "FAQs" },
		{ href: "/#support", text: "Support The Maintainers" },
	];
</script>

<nav
	class={cx([
		"sticky left-0 top-0 z-[49] translate-y-0 transition-all",
		"bg-base-100/80 backdrop-blur md:flex md:items-center md:py-5",
	])}
>
	<div
		class="flex h-16 items-center justify-between px-5 md:inline md:h-full md:justify-start"
	>
		<img
			src={Logo}
			alt="logo"
		/>
		<button
			on:click={toggleOpen}
			class="rounded-md bg-inherit p-2 transition-all hover:bg-[#3a3a3a] md:hidden"
		>
			<Hamburger />
		</button>
	</div>

	<div
		class="absolute left-0 top-16 w-full pr-4 md:static md:ml-5 md:mt-1 md:w-max"
	>
		<ul
			class={cx([
				"bg-base-100/80 h-0 w-full py-0 backdrop-blur",
				{ "h-[350px] py-5": show },
				"flex flex-col items-center justify-between",
				"transition-all",
				"relative overflow-clip",
				"z-50 gap-10 md:h-full md:flex-row md:bg-inherit",
			])}
		>
			<li class="md:hidden">
				<Dashboard />
			</li>
			{#each navMap as { href, text } (text)}
				<li>
					<a
						class="text-neutral-content font-bold hover:underline"
						{href}
					>
						{text}
					</a>
				</li>
			{/each}
		</ul>
	</div>

	<div class="ml-auto hidden gap-5 px-5 md:flex">
		<Dashboard />
	</div>
</nav>
