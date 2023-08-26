<script lang="ts">
	import cx from "classnames";

	export let className: string = "";
</script>

<div class={cx(["row border-light/5 w-full border-b px-4 py-2", className])}>
	<slot />
</div>

<style lang="postcss">
	.row {
		display: grid;
		grid-template-columns: 75px 15ch 1fr 15ch 15ch;
		grid-gap: 10px;
	}
</style>
