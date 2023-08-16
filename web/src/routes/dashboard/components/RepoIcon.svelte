<script lang="ts">
  export let tooltip: string
  export let show: boolean
</script>

{#if show}
  <div class="tooltip tooltip-primary" data-tip={tooltip}>
    <slot />
  </div>
{/if}
