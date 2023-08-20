<script lang="ts">
</script>

<div class="row w-full border-b border-light/5 py-2 px-4">
	<slot />
</div>

<style lang="postcss">
	.row {
		display: grid;
		grid-template-columns: 85px 25ch 1fr 15ch 15ch;
		grid-gap: 10px;
	}
</style>
