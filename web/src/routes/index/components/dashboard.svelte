<script lang="ts">
	import { AuthStore, oauth, type User } from "@lib/auth";
	import { onMount } from "svelte";

	let user: User | undefined;
	onMount(() => {
		user = AuthStore.user();
	});

	let isLoggedIn: boolean | undefined = undefined;
	$: isLoggedIn = user !== undefined;
</script>

{#if isLoggedIn !== undefined}
	{#if isLoggedIn}
		<a
			href="/dashboard"
			class="font-bold"
		>
			Dashboard
		</a>
	{:else}
		<button
			on:click={() => oauth("/dashboard")}
			class="md:btn md:btn-primary font-bold md:normal-case"
		>
			Sign in <span>
				<i class="fa-brands fa-github-alt"></i>
			</span>
		</button>
	{/if}
{/if}
