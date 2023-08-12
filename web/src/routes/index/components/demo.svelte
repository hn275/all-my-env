<script lang="ts">
	export let encrypted: boolean;

	type Variables = {
		key: string;
		value: string;
		encoded: string;
	};
	const demoVar: Variables[] = [
		{
			key: "API_KEY",
			value: "api-secret-key",
			encoded: "Ym05M1lYeDV",
		},
		{
			key: "TOKEN",
			value: "my-token",
			encoded: "YXVCZ29ZQm02ZGR2",
		},
		{
			key: "PG_DSN",
			value: "pg://foo:bar@153.531.6.2",
			encoded: "cGc6L2ZvbzpiY",
		},
		{
			key: "AWS_KEY",
			value: "bucket",
			encoded: "c2Q3YWZnZDkwZjg5YThk",
		},
		{
			key: "AWS_SECRET",
			value: "password",
			encoded: "Zjd0c2RmNzk4YTd",
		},
		{
			key: "JWT_SECRET",
			value: "super-secure-secret",
			encoded: "OGZkN2E2c2=",
		},
	];
</script>

<ul
	class="bg-dark-100 flex h-full w-full flex-col gap-2 rounded-lg p-5 lg:w-96"
>
	{#each demoVar as { encoded, key, value }, i}
		<li class="text-main flex w-96 flex-row font-semibold">
			<p class="font-normal text-white">{key}=</p>
			<p>
				{#if encrypted}
					"{encoded}"
				{:else}
					"{value}"
				{/if}
			</p>
		</li>
	{/each}
</ul>
