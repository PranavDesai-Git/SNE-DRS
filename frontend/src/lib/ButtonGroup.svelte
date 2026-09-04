<script>
  export let options = []; // Array of { label: string, value: string }
  export let selected = options[0]?.value;
</script>

<div class="gov-btn-group">
  {#each options as option}
    <button 
      class="group-btn" 
      class:active={selected === option.value}
      on:click={() => selected = option.value}
    >
      {option.label}
    </button>
  {/each}
</div>

<style>
  .gov-btn-group {
    display: inline-flex;
    border: 1px solid var(--secondary);
    border-radius: 4px;
    overflow: hidden;
  }

  .group-btn {
    background-color: var(--background);
    border: none;
    border-right: 1px solid var(--secondary);
    padding: 0.5rem 1rem;
    font-family: 'Poppins', sans-serif;
    font-size: 0.9rem;
    color: var(--text);
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .group-btn:last-child {
    border-right: none;
  }

  .group-btn:hover:not(.active) {
    background-color: color-mix(in srgb, var(--background), black 4%);
  }

  .group-btn.active {
    background-color: var(--primary-dark);
    color: var(--background);
  }
</style>
