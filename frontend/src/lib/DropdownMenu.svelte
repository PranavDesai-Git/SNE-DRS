<script>
  export let label = "Actions";
  export let items = []; // Array of { text: string, action: function }
  
  let open = false;

  function toggle() {
    open = !open;
  }

  function handleAction(action) {
    action();
    open = false;
  }
</script>

<div class="gov-dropdown-container">
  <button class="dropdown-trigger" on:click={toggle}>
    {label}
    <span class="caret">▼</span>
  </button>
  
  {#if open}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="dropdown-overlay" on:click={() => open = false}></div>
    <div class="dropdown-content">
      {#each items as item}
        <button class="dropdown-item" on:click={() => handleAction(item.action)}>
          {item.text}
        </button>
      {/each}
    </div>
  {/if}
</div>

<style>
  .gov-dropdown-container {
    position: relative;
    display: inline-block;
  }

  .dropdown-trigger {
    background-color: var(--background);
    border: 1px solid var(--secondary);
    color: var(--text);
    padding: 0.4rem 0.75rem;
    border-radius: 4px;
    font-family: 'Poppins', sans-serif;
    font-size: 0.85rem;
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    transition: border-color 0.2s;
  }

  .dropdown-trigger:hover {
    border-color: var(--primary-dark);
  }

  .caret {
    font-size: 0.6rem;
  }

  .dropdown-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    z-index: 999;
  }

  .dropdown-content {
    position: absolute;
    top: calc(100% + 4px);
    right: 0;
    background-color: var(--background);
    border: 1px solid var(--secondary);
    border-radius: 4px;
    min-width: 150px;
    z-index: 1000;
    display: flex;
    flex-direction: column;
    box-shadow: inset 0 2px 6px rgba(0, 0, 0, 0.08); /* Using inset shadow to comply with no-float rule */
  }

  .dropdown-item {
    background: none;
    border: none;
    padding: 0.5rem 1rem;
    text-align: left;
    font-family: 'Poppins', sans-serif;
    font-size: 0.85rem;
    color: var(--text);
    cursor: pointer;
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  }

  .dropdown-item:last-child {
    border-bottom: none;
  }

  .dropdown-item:hover {
    background-color: color-mix(in srgb, var(--background), black 4%);
  }
</style>
