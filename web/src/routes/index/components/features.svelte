<script lang="ts">
	import cx from "classnames";
	import ShareSVG from "../assets/share.svg";
	import Demo from "./demo.svelte";

	type FeatureType = {
		title: string;
		content: string;
	};
	const features: FeatureType[] = [
		{
			title: "First Thing",
			content:
				"Sit velit molestiae nesciunt aperiam ex. Fuga eius sapiente veritatis perspiciatis officia! Deserunt ut exercitationem accusamus quis iure. Deserunt ducimus ab rem voluptates ad cum, obcaecati consequuntur. Magnam nobis officiis.",
		},
		{
			title: "Second Thing",
			content:
				"Dolor dolorem officiis sed reprehenderit eveniet Iusto doloremque hic sint culpa nam illo ipsa minus culpa Cum eligendi atque delectus",
		},
		{
			title: "Third First",
			content:
				"Lorem cumque veniam tenetur corrupti molestiae? In necessitatibus sequi dolore neque sed? Aut explicabo voluptates laudantium culpa ratione? Fuga corrupti!",
		},
	];

	const pages: number[] = [];
	for (let i = 0; i < features.length; ++i) pages.push(i);
	let activeTile: number = 0;
</script>

<section
	class={cx([
		"mx-auto flex w-full flex-col items-center justify-center",
		"gap-8 lg:flex-row lg:gap-24",
	])}
>
	<div class="mt-5 flex w-3/4 lg:w-1/2 flex-col items-start">
		<div>
			<div class="h-64 w-full lg:h-48 lg:w-[40ch]">
				<h3 class="font-bold text-xl mb-5">
					{features[activeTile].title}
				</h3>
				<p>{features[activeTile].content}</p>
			</div>
		</div>

		<ul class="mt-3 flex items-center justify-start gap-3">
			{#each pages as page}
				<li>
					<button
						class={cx([
							"h-8 w-8 rounded-lg border transition-all",
							{
								"bg-dark-100 border-dark-100":
									page !== activeTile,
								"hover:border-fuchsia-400 hover:text-fuchsia-400/50":
									page !== activeTile,
                "border-fuchsia-400/50 bg-fuchsia-400/10 text-fuchsia-400":
									page === activeTile,
							},
						])}
						on:click={() => (activeTile = page)}>{page}</button
					>
				</li>
			{/each}
		</ul>
	</div>

	<div class="w-[35ch] lg:w-1/2">
		<div>
			{#if activeTile === 2}
				<div class="w-full">
					<img
						src={ShareSVG}
						role="presentation"
						alt=""
						class="max-h-52"
					/>
				</div>
			{:else}
				<div class="h-full">
					<Demo encrypted={activeTile === 1} />
				</div>
			{/if}
		</div>
	</div>
</section>
