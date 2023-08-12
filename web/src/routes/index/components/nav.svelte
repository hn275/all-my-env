<script lang="ts">
	import cx from "classnames";
	import Logo from "@assets/logo.svg";
	import StarUs from "./star.svelte";
	import LogInBtn from "./login.svelte";
	import Hamburger from "../assets/hamburger.svelte";

	export let loading: boolean;

	let show: boolean = false;
	const toggleOpen: () => void = () => (show = !show);
</script>

<nav
	class={cx([
		"sticky left-0 top-0 z-[49] transition-all translate-y-0",
		"bg-dark justify-between bg-dark-200/60 backdrop-blur md:flex",
	])}
>
	<div class="flex h-16 items-center justify-between px-5">
		<img src={Logo} alt="logo" />
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
				`w-full h-0 py-0`,
				{ "h-[350px] py-5": show },
				"text-light bg-[#1e1e1e] font-semibold",
				"flex flex-col items-center justify-between",
				"transition-all",
				"relative overflow-clip",
				"z-50 gap-10 md:h-full md:flex-row md:bg-inherit",
			])}
		>
			<li>
				<LogInBtn {loading} />
			</li>
			<li>
				<a href="/pricing">Pricing</a>
			</li>
			<li>
				<a href="/docs">Docs</a>
			</li>
			<li>
				<a href="/faq">FAQ</a>
			</li>
			<li>
				<StarUs />
			</li>
		</ul>
	</div>
</nav>
