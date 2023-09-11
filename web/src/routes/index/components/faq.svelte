<script lang="ts">
	export let header: string;
</script>

<div class="join-item border border-neutral collapse collapse-arrow">
	<input type="checkbox" />
	<h5 class="collapse-title font-bold text-lg">{header}</h5>
	<div class="collapse-content">
		<slot />
	</div>
</div>
