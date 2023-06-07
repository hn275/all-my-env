<script lang="ts">
	import { onMount } from "svelte";
	import { Cipher, SecretKeys } from "@mod/lib/cipher";
	import { API } from "@mod/lib/routes";

	let code = "";
	let encryptedStr: string = "";
	let decryptedStr: string = "";
	let ivStr: string = "";
	let status = "";

	onMount(() => {
		code = new URLSearchParams(window.location.search).get(
			"code",
		) as string;

		(async () => {
			const url = `${API}/auth`;
			console.log(url);
			const response = await fetch(url, {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({ code }),
			});
			status = response.statusText;
		})();
	});
</script>

<p>verifying</p>
<p>{code}</p>
<p>{status}</p>
<p>{decryptedStr === code}</p>
