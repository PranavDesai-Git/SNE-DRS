<script>
  export let label = "";
  export let value = "";
  export let placeholder = "Search...";

  function clearSearch() {
    value = "";
  }
</script>

<div class="input-group">
  {#if label}
    <label>{label}</label>
  {/if}
  <div class="search-wrapper">
    <!-- SVG Search Icon -->
    <svg class="search-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="11" cy="11" r="8"></circle>
      <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
    </svg>
    
    <input type="text" bind:value {placeholder} class="gov-search" />
    
    {#if value}
      <button class="clear-btn" on:click={clearSearch}>&times;</button>
    {/if}
  </div>
</div>

<style>
  .input-group {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    width: 100%;
  }
  
  label {
    font-size: 0.85rem;
    font-weight: 400;
  }

  .search-wrapper {
    position: relative;
    display: flex;
    align-items: center;
  }

  .search-icon {
    position: absolute;
    left: 0.75rem;
    color: color-mix(in srgb, var(--text) 60%, transparent);
    pointer-events: none;
  }

  .gov-search {
    font-family: 'Poppins', sans-serif;
    font-size: 1rem;
    width: 100%;
    padding: 0.5rem 2rem 0.5rem 2.25rem;
    border: 1px solid var(--secondary);
    border-radius: 4px;
    background-color: var(--background);
    color: var(--text);
    outline: none;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
  }

  .gov-search:focus {
    border-color: var(--primary-dark);
    box-shadow: 0 0 0 2px rgba(42, 71, 36, 0.2);
  }

  .clear-btn {
    position: absolute;
    right: 0.5rem;
    background: none;
    border: none;
    font-size: 1.25rem;
    cursor: pointer;
    color: color-mix(in srgb, var(--text) 60%, transparent);
    line-height: 1;
    padding: 0 0.25rem;
  }

  .clear-btn:hover {
    color: var(--text);
  }
</style>
