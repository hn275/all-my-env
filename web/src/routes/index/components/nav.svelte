<script lang="ts">
	import cx from "classnames";
	import Logo from "@assets/logo.svg";
	import StarUs from "./star.svelte";
	import LogInBtn from "./login.svelte";
	import Hamburger from "../assets/hamburger.svelte";

	export let loading: boolean;

	let show: boolean = false;
	const toggleOpen: () => void = () => (show = !show);

	type NavMap = {
		href: string;
		text: string;
	};
	const navMap: Array<NavMap> = [
		{ href: "/pricing", text: "Pricing" },
		{ href: "/docs", text: "Docs" },
		{ href: "/faq", text: "FAQ" },
	];
</script>

<nav
	class={cx([
		"sticky left-0 top-0 z-[49] translate-y-0 transition-all",
		"bg-dark justify-between backdrop-blur md:flex",
	])}
>
	<div class="flex h-16 items-center justify-between px-5">
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

	<div class="absolute left-0 top-16 w-full pr-4 md:static md:w-max">
		<ul
			class={cx([
				"h-0 w-full py-0",
				{ "h-[350px] py-5": show },
				"text-light md:text-light/60 bg-[#1e1e1e]",
				"flex flex-col items-center justify-between",
				"transition-all",
				"relative overflow-clip",
				"z-50 gap-10 md:h-full md:flex-row md:bg-inherit",
			])}
		>
			<li>
				<StarUs />
			</li>
			{#each navMap as { href, text } (text)}
				<li>
					<a
						class="transition hover:text-light"
						{href}
					>
						{text}
					</a>
				</li>
			{/each}
			<li>
				<LogInBtn {loading} />
			</li>
		</ul>
	</div>
</nav>
