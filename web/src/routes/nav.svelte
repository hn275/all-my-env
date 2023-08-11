<script lang="ts">
	import cx from "classnames";
	import Logo from "@assets/logo.svg";
	import StarUs from "./star.svelte";
	import LogInBtn from "./login.svelte";
	import { oauth } from "@lib/auth";

	export let loading: boolean;

	let show: boolean = true;
	function toggleOpen(): void {
		show != show;
	}
</script>

<nav
	class={cx([
		"sticky left-0 top-0 z-[49] transition-all",
		"bg-dark -translate-y-full justify-between bg-transparent backdrop-blur md:flex",
		{ "translate-y-0": show },
	])}
>
	<div class="flex h-16 items-center justify-between px-5">
		<img src={Logo} alt="logo" />
		<button
			on:click={toggleOpen}
			class="rounded-md bg-inherit p-2 transition-all hover:bg-[#3a3a3a] md:hidden"
		>
			menu button
		</button>
	</div>

	<div class={cx(["absolute left-0 top-16 w-full", "md:static md:w-max"])}>
		<ul
			class={cx([
				`w-full ${show ? "h-[350px] py-5" : "h-0 py-0"}`,
				"text-light bg-[#1e1e1e] font-semibold",
				"flex flex-col items-center justify-between",
				"transition-all",
				"relative overflow-clip",
				"z-50 gap-10 md:h-full md:flex-row md:bg-inherit",
			])}
		>
			<li>
				<LogInBtn
					{loading}
					handleClick={() => oauth()}
				/>
			</li>
			<li>
				<a href="/pricing">Pricing</a>
			</li>
			<li>
				<a href="docs">Docs</a>
			</li>
			<li>
				<a href="faq">FAQ</a>
			</li>
			<li>
				<StarUs />
			</li>
			<li>
				<button
					on:click={toggleOpen}
					class="absolute bottom-2 left-1/2 -translate-x-1/2 md:hidden"
				>
					close menu button
				</button>
			</li>
		</ul>
	</div>
</nav>
