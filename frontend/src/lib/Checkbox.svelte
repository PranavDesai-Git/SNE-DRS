<script>
  export let label = "";
  export let checked = false;
</script>

<label class="checkbox-container">
  <input type="checkbox" bind:checked />
  <span class="checkmark"></span>
  {#if label}
    <span class="label-text">{label}</span>
  {/if}
</label>

<style>
  .checkbox-container {
    display: flex;
    align-items: center;
    position: relative;
    cursor: pointer;
    font-size: 0.85rem;
    font-family: 'Poppins', sans-serif;
    user-select: none;
  }

  .checkbox-container input {
    position: absolute;
    opacity: 0;
    cursor: pointer;
    height: 0;
    width: 0;
  }

  .checkmark {
    height: 18px;
    width: 18px;
    background-color: var(--background);
    border: 1px solid var(--secondary);
    border-radius: 4px;
    margin-right: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .checkbox-container:hover .checkmark {
    border-color: var(--primary-dark);
    transform: scale(1.05);
  }

  .checkbox-container:active .checkmark {
    transform: scale(0.95);
  }

  .checkbox-container input:checked ~ .checkmark {
    background-color: var(--primary-dark);
    border-color: var(--primary-dark);
  }

  .checkmark:after {
    content: "";
    display: none;
    width: 4px;
    height: 8px;
    border: solid #ffffff;
    border-width: 0 2px 2px 0;
    transform: rotate(45deg);
    margin-bottom: 2px;
  }

  .checkbox-container input:checked ~ .checkmark:after {
    display: block;
  }
</style>
