<script>
  export let message = "";
  export let type = "success"; // "success", "error", "info"
  export let visible = false;
</script>

{#if visible}
  <div class="notification {type}">
    <div class="message">{message}</div>
    <button class="close-btn" on:click={() => visible = false}>&times;</button>
  </div>
{/if}

<style>
  .notification {
    position: fixed;
    top: 2rem;
    left: 50%;
    transform: translateX(-50%);
    padding: 1rem 2rem;
    border-radius: 4px;
    background-color: var(--background);
    color: var(--text);
    box-shadow: 0 4px 15px rgba(0,0,0,0.15);
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 3rem;
    z-index: 10000;
    border: 1px solid var(--secondary);
    font-family: 'Poppins', sans-serif;
    font-size: 0.9rem;
    animation: slideDown 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: var(--text);
    opacity: 0.5;
    transition: opacity 0.2s;
    padding: 0;
    line-height: 1;
  }

  .close-btn:hover {
    opacity: 1;
  }

  @keyframes slideDown {
    from { transform: translate(-50%, -120%); opacity: 0; }
    to { transform: translate(-50%, 0); opacity: 1; }
  }
</style>
