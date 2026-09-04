<script>
  export let isOpen = false;
  
  function close() {
    isOpen = false;
  }
</script>

{#if isOpen}
  <div class="gov-popover">
    <button class="close-btn" on:click={close}>&times;</button>
    <div class="popover-content">
      <slot />
    </div>
  </div>
{/if}

<style>
  .gov-popover {
    position: absolute;
    background-color: var(--background);
    border: 1px solid var(--secondary);
    border-radius: 4px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);
    z-index: 1000;
    min-width: 200px;
  }

  .close-btn {
    position: absolute;
    top: 0.25rem;
    right: 0.25rem;
    background: none;
    border: none;
    font-size: 1.25rem;
    cursor: pointer;
    color: color-mix(in srgb, var(--text) 60%, transparent);
    line-height: 1;
    padding: 0 0.25rem;
  }

  .close-btn:hover {
    color: var(--text);
  }

  .popover-content {
    padding: 1rem;
    padding-top: 1.5rem;
  }
</style>
