<script>
  export let value = 0; // 0 to 100
  export let label = "";
</script>

<div class="progress-container">
  {#if label}
    <div class="progress-label">
      <span>{label}</span>
      <span>{value}%</span>
    </div>
  {/if}
  <div class="progress-track">
    <div 
      class="progress-fill" 
      style="width: {Math.min(Math.max(value, 0), 100)}%"
    ></div>
  </div>
</div>

<style>
  .progress-container {
    width: 100%;
    margin-bottom: 1rem;
  }

  .progress-label {
    display: flex;
    justify-content: space-between;
    font-family: 'Poppins', sans-serif;
    font-size: 0.85rem;
    color: var(--text);
    margin-bottom: 0.25rem;
  }

  .progress-track {
    width: 100%;
    height: 8px;
    background-color: var(--secondary);
    border-radius: 999px;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background-color: var(--primary-dark);
    transition: width 0.3s ease;
  }
</style>
