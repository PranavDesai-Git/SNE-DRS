<script>
  export let label = "";
  export let min = 0;
  export let max = 100;
  export let value = 50;
</script>

<div class="input-group">
  {#if label}
    <div class="label-row">
      <label>{label}</label>
      <span class="value-display">{value}</span>
    </div>
  {/if}
  <input type="range" {min} {max} bind:value class="gov-slider" />
</div>

<style>
  .input-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }
  
  .label-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  label, .value-display {
    font-size: 0.85rem;
    font-weight: 400;
  }

  .gov-slider {
    -webkit-appearance: none;
    appearance: none;
    width: 100%;
    height: 6px;
    background: var(--secondary);
    border-radius: 999px;
    outline: none;
  }

  .gov-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: var(--color-habitation); /* Used a darker green from the theme for contrast */
    cursor: pointer;
    box-shadow: 0 2px 4px rgba(0,0,0,0.2);
    transition: transform 0.1s ease;
  }

  .gov-slider::-webkit-slider-thumb:hover {
    transform: scale(1.15);
  }
</style>
