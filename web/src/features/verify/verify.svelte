<script lang="ts">
	import { onMount } from "svelte";
	import { API } from "@mod/lib/routes";
	import type { GithubAccount } from "@mod/schemas/github";

	let code = "";
	let status = "";
	let account: GithubAccount;

	onMount(() => {
		code = new URLSearchParams(window.location.search).get(
			"code",
		) as string;

		(async () => {
			const url = `${API}/auth`;
			const response = await fetch(url, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ code }),
			});
			status = response.statusText;
			account = await response.json();
		})();
	});
</script>

<p>verifying</p>
<p>{code}</p>
<p>{status}</p>
<p>{JSON.stringify(account)}</p>
