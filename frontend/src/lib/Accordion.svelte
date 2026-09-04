<script>
  export let title = "Expand Details";
  export let open = false;
</script>

<div class="accordion">
  <button class="accordion-header" on:click={() => open = !open}>
    <span>{title}</span>
    <span class="chevron">{open ? '−' : '+'}</span>
  </button>
  {#if open}
    <div class="accordion-content">
      <slot />
    </div>
  {/if}
</div>

<style>
  .accordion {
    border: 1px solid var(--secondary);
    border-radius: 4px;
    margin-bottom: 0.5rem;
    overflow: hidden;
  }

  .accordion-header {
    width: 100%;
    background-color: var(--background);
    border: none;
    padding: 0.75rem 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    font-family: 'Poppins', sans-serif;
    font-size: 0.9rem;
    color: var(--text);
    transition: background-color 0.2s;
  }

  .accordion-header:hover {
    background-color: color-mix(in srgb, var(--background), var(--secondary) 15%);
  }

  .chevron {
    font-size: 1.2rem;
    line-height: 1;
    color: var(--primary-dark);
  }

  .accordion-content {
    padding: 1rem;
    border-top: 1px solid var(--secondary);
    background-color: color-mix(in srgb, var(--background), black 2%);
  }
</style>
