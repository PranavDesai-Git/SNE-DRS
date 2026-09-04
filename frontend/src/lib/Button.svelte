<script>
  // Exports allow us to pass properties to the component
  export let type = "button";
  export let disabled = false;
  export let fullWidth = false;
</script>

<!-- 
  We use <slot /> so you can put text or icons inside the button when you use it.
  The on:click forwards the click event to whatever parent uses this component.
-->
<button 
  {type}
  {disabled}
  class="gov-button"
  class:full-width={fullWidth}
  on:click
>
  <slot />
</button>

<style>
  .gov-button {
    font-family: 'Poppins', sans-serif;
    font-size: 1rem;
    padding: 0.5rem 1rem;
    border: 1px solid transparent;
    border-radius: 4px; /* Slight rounding for a professional look */
    background-color: var(--primary);
    color: var(--text); /* Use the global text color (#0f1510) */
    cursor: pointer;
    transition: transform 0.15s ease, box-shadow 0.15s ease;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }

  .gov-button:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  }

  .gov-button:active:not(:disabled) {
    transform: translateY(1px);
    box-shadow: none;
  }

  .gov-button:disabled {
    background-color: var(--secondary);
    color: var(--text);
    opacity: 0.6;
    cursor: not-allowed;
    box-shadow: none;
  }

  .full-width {
    width: 100%;
  }
</style>
