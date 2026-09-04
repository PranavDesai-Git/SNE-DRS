<script>
  export let options = []; // Array of { value: string, label: string }
  export let selected = "";
  export let name = "radio-group";
</script>

<div class="gov-radio-group">
  {#each options as option}
    <label class="radio-label">
      <input 
        type="radio" 
        {name} 
        value={option.value} 
        bind:group={selected} 
        class="gov-radio"
      />
      <span class="radio-text">{option.label}</span>
    </label>
  {/each}
</div>

<style>
  .gov-radio-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .radio-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    font-size: 0.9rem;
    color: var(--text);
  }

  .gov-radio {
    accent-color: var(--primary-dark);
    width: 1rem;
    height: 1rem;
    cursor: pointer;
  }
</style>
