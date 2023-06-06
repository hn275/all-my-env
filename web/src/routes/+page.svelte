<script lang="ts">
	import Html from "@mod/layout/html.svelte";

	const params = new URLSearchParams({
		client_id: "567d9b90135f2031eb22",
		scope: "repo:read user:read",
	});
	const oauth = new URL("/login/oauth/authorize", "https://github.com");
	oauth.search = params.toString();
</script>

<Html title="login">
	<section class="flex flex-col gap-3 justify-center items-center">
		<h1 class="text-2xl m-3">Company name or something</h1>
		<a
			class="text-lg text-white bg-blue-300 p-3 rounded-md"
			href={oauth.toString()}
			target="_blank"
		>
			Sign in with Github
		</a>
	</section>
</Html>
