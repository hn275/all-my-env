<script lang="ts">
	import { onMount } from "svelte";
	import {
		PUBLIC_GITHUB_CLIENT_ID,
		PUBLIC_NODE_ENV,
	} from "$env/static/public";

	let oauth: string = "/dash";
	onMount(() => {
		const isProd: boolean = PUBLIC_NODE_ENV === "production";
		const http = isProd ? "https://" : "http://";
		const redirect_uri = http + window.location.host + "/auth";
		const client_id = PUBLIC_GITHUB_CLIENT_ID;
		const scope = "repo read:user read:org";
		const p = new URLSearchParams({ client_id, redirect_uri, scope });
		oauth = "https://github.com/login/oauth/authorize?" + p.toString();
	});
</script>

<a
	href={oauth}
	class="btn btn-primary"
>
	Sign in with GitHub
	<span>
		<i class="fa-brands fa-github fa-xl"></i>
	</span>
</a>
