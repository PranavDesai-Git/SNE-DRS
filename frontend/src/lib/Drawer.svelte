<script>
  export let isOpen = false;
  export let title = "";
  export let position = "right"; // "left" or "right"

  function closeDrawer() {
    isOpen = false;
  }
</script>

{#if isOpen}
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <div class="drawer-backdrop" on:click={closeDrawer}></div>
{/if}

<div class="gov-drawer {position}" class:open={isOpen}>
  <div class="drawer-header">
    {#if title}
      <h3>{title}</h3>
    {/if}
    <button class="close-btn" on:click={closeDrawer}>&times;</button>
  </div>
  <div class="drawer-content">
    <slot />
  </div>
</div>

<style>
  .drawer-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.4);
    z-index: 9998;
  }

  .gov-drawer {
    position: fixed;
    top: 0;
    height: 100vh;
    width: 400px;
    max-width: 90vw;
    background-color: var(--background);
    z-index: 9999;
    transition: transform 0.3s cubic-bezier(0.165, 0.84, 0.44, 1);
    box-shadow: 0 0 15px rgba(0,0,0,0.1);
    display: flex;
    flex-direction: column;
  }

  .right {
    right: 0;
    transform: translateX(100%);
  }

  .left {
    left: 0;
    transform: translateX(-100%);
  }

  .open {
    transform: translateX(0);
  }

  .drawer-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid var(--secondary);
  }

  .drawer-header h3 {
    margin: 0;
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: var(--text);
    padding: 0;
    line-height: 1;
    opacity: 0.6;
  }

  .close-btn:hover {
    opacity: 1;
  }

  .drawer-content {
    padding: 1.5rem;
    flex: 1;
    overflow-y: auto;
  }
</style>
