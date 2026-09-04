<script>
  export let text = "";
</script>

<span class="gov-tag">
  {#if text}
    {text}
  {:else}
    <slot />
  {/if}
</span>

<style>
  .gov-tag {
    display: inline-flex;
    align-items: center;
    padding: 0.15rem 0.5rem;
    border: 1px solid var(--secondary);
    border-radius: 4px;
    font-size: 0.75rem;
    color: var(--text);
    background-color: transparent;
  }
</style>
