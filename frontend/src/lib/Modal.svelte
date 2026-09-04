<script>
  export let isOpen = false;
  export let title = "";
</script>

{#if isOpen}
  <div class="modal-backdrop" on:click={() => isOpen = false}>
    <div class="modal-content" on:click|stopPropagation>
      {#if title}
        <div class="modal-header">
          <h3>{title}</h3>
          <button class="close-btn" on:click={() => isOpen = false}>&times;</button>
        </div>
      {/if}
      <div class="modal-body">
        <slot />
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    top: 0; left: 0; width: 100vw; height: 100vh;
    background-color: rgba(0, 0, 0, 0.4);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
    backdrop-filter: blur(2px);
  }

  .modal-content {
    background-color: var(--background);
    border: 1px solid var(--secondary);
    border-radius: 4px;
    width: 90%;
    max-width: 500px;
    box-shadow: 0 10px 25px rgba(0,0,0,0.2);
    animation: slideUp 0.2s ease-out;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    border-bottom: 1px solid var(--secondary);
  }

  .modal-header h3 {
    margin: 0;
    font-size: 1.2rem; 
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: var(--text);
  }

  .modal-body {
    padding: 1.5rem;
  }

  @keyframes slideUp {
    from { transform: translateY(20px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
  }
</style>
