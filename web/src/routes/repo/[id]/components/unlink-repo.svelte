<script lang="ts">
	import cx from "classnames";
</script>

<button
	class={cx([
		"btn btn-outline border-error/20 text-error/50 hover:text-error",
		"hover:border-error/80 hover:bg-dark-200 text-xs",
	])}
>
	<i class="fa-solid fa-link-slash" />
	Delete repository
</button>
