<script>
  export let label = "";
  export let value = "";
  export let placeholder = "";
  export let rows = 4;
</script>

<div class="input-group">
  {#if label}
    <label>{label}</label>
  {/if}
  <textarea bind:value {placeholder} {rows} class="gov-textarea"></textarea>
</div>

<style>
  .input-group {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  label {
    font-size: 0.85rem;
    font-weight: 400;
  }

  .gov-textarea {
    font-family: 'Poppins', sans-serif;
    font-size: 1rem;
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--secondary);
    border-radius: 4px;
    background-color: var(--background);
    color: var(--text);
    outline: none;
    resize: vertical;
    transition: border-color 0.2s ease, box-shadow 0.2s ease;
  }

  .gov-textarea:focus {
    border-color: var(--primary-dark);
    box-shadow: 0 0 0 2px rgba(42, 71, 36, 0.2);
  }
</style>
