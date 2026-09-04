<script>
  export let title = "";
</script>

<div class="gov-card">
  {#if title}
    <div class="card-header">
      <h4>{title}</h4>
    </div>
  {/if}
  
  <slot name="band" />

  <div class="card-body">
    <slot />
  </div>
</div>

<style>
  .gov-card {
    background-color: rgba(255, 255, 255, 0.3);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    border: 1px solid rgba(255, 255, 255, 0.6);
    border-radius: 4px;
    overflow: hidden;
    box-shadow: inset 0 2px 6px rgba(0, 0, 0, 0.05); /* Inward shadow */
  }

  .card-header {
    padding: 1rem 1rem 0.5rem 1rem;
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  }

  .card-header h4 {
    margin: 0;
    color: var(--text);
  }

  .card-body {
    padding: 1rem;
  }
</style>
