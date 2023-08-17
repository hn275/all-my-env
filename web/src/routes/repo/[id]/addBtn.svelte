<script lang="ts">
	export let writeAccess: boolean;
</script>

<button
	class="btn btn-primary text-xs"
	disabled={!writeAccess}
	on:click
>
	<i class="fa-solid fa-plus" />
	Add
</button>
