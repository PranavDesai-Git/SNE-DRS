<script>
  export let label = "";
  export let value = "";
  export let options = []; // Array of { value, label }
</script>

<div class="input-group">
  {#if label}
    <label>{label}</label>
  {/if}
  <select bind:value class="gov-select">
    {#each options as option}
      <option value={option.value}>{option.label}</option>
    {/each}
  </select>
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

  .gov-select {
    font-family: 'Poppins', sans-serif;
    font-size: 1rem;
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--secondary);
    border-radius: 4px;
    background-color: var(--background);
    color: var(--text);
    outline: none;
    transition: border-color 0.2s;
    cursor: pointer;
    appearance: auto; /* Uses native dropdown arrow but styles the box */
  }

  .gov-select:focus {
    border-color: var(--primary-dark);
  }
</style>
