<script>
  // Static legend mapping based on the design system
  const items = [
    { label: 'Immediate Risk', class: 'immediate' },
    { label: 'Moderate Risk', class: 'moderate' },
    { label: 'Safe Zone', class: 'safe' },
    { label: 'Habitation', class: 'habitation' }
  ];
</script>

<div class="gov-map-legend">
  <h4 class="legend-title">Hazard Zones</h4>
  <div class="legend-items">
    {#each items as item}
      <div class="legend-item">
        <span class="color-box {item.class}"></span>
        <span class="legend-label">{item.label}</span>
      </div>
    {/each}
  </div>
</div>

<style>
  .gov-map-legend {
    background-color: var(--background);
    border: 1px solid var(--secondary);
    border-radius: 4px;
    padding: 1rem;
    display: inline-block;
  }

  .legend-title {
    margin-top: 0;
    margin-bottom: 0.75rem;
    font-size: 0.95rem;
  }

  .legend-items {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .color-box {
    width: 16px;
    height: 16px;
    border-radius: 2px;
    display: inline-block;
  }

  .legend-label {
    font-size: 0.85rem;
    color: var(--text);
  }

  /* Exact colors mapped from app.css */
  .immediate { background-color: var(--color-redzone); }
  .moderate { background-color: var(--color-moderate); }
  .safe { background-color: var(--color-safe); }
  .habitation { background-color: var(--color-habitation); }
</style>
