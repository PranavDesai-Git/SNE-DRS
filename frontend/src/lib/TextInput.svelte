<script>
  export let label = "";
  export let value = "";
  export let placeholder = "";
</script>

<div class="input-group">
  {#if label}
    <label>{label}</label>
  {/if}
  <input type="text" bind:value {placeholder} class="gov-input" />
</div>

<style>
  .input-group {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    margin-bottom: 1rem;
  }
  
  label {
    font-size: 0.85rem;
    font-weight: 400;
  }

  .gov-input {
    font-family: 'Poppins', sans-serif;
    font-size: 1rem;
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--secondary);
    border-radius: 4px;
    background-color: var(--background);
    color: var(--text);
    outline: none;
    transition: border-color 0.2s;
  }

  .gov-input:focus {
    border-color: var(--primary-dark);
    box-shadow: 0 0 0 2px rgba(42, 71, 36, 0.2); /* Soft primary-dark glow */
  }
</style>
