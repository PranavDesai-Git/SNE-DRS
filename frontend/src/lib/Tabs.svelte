<script>
  export let tabs = []; // array of strings
  export let activeTab = tabs[0];
</script>

<div class="tabs-container">
  {#each tabs as tab}
    <button 
      class="tab" 
      class:active={activeTab === tab} 
      on:click={() => activeTab = tab}
    >
      {tab}
    </button>
  {/each}
</div>

<style>
  .tabs-container {
    display: flex;
    border-bottom: 2px solid var(--secondary);
    margin-bottom: 1rem;
    gap: 1rem;
  }

  .tab {
    background: none;
    border: none;
    font-family: 'Poppins', sans-serif;
    font-size: 0.85rem;
    color: color-mix(in srgb, var(--text) 60%, transparent);
    padding: 0.5rem 0.25rem;
    cursor: pointer;
    position: relative;
    transition: color 0.2s;
  }

  .tab:hover {
    color: var(--text);
  }

  .tab.active {
    color: var(--text);
  }

  .tab.active::after {
    content: '';
    position: absolute;
    bottom: -2px; /* cover the container border */
    left: 0;
    width: 100%;
    height: 2px;
    background-color: var(--primary-dark);
  }
</style>
