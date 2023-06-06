<script lang="ts">
	import Html from "@mod/layout/html.svelte";
	import { onMount } from "svelte";
	import { Cipher, SecretKeys } from "@mod/lib/cipher"; // TODO: implement this

	let code = "";

	let encrypted: string = "";
	let decrypted: string = "";

	onMount(() => {
		code = new URLSearchParams(window.location.search).get("code") as string;
		const aesEncryptor = new Cipher(code);
		aesEncryptor.setKey(SecretKeys.Auth);
		encrypted = aesEncryptor.encrypt();

		const aesDecryptor = new Cipher(encrypted);
		aesDecryptor.setKey(SecretKeys.Auth);
		decrypted = aesDecryptor.decrypt();
	});
</script>

<Html title="verify">
	<p>verifying</p>
	<p>{code}</p>
	<p>encrypted {encrypted}</p>
	<p>decrypted {decrypted}</p>
	<p>{decrypted === code}</p>
</Html>
