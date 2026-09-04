<script>
  export let label = "";
  export let value = "0";
</script>

<div class="metric">
  <span class="metric-label">{label}</span>
  <span class="metric-value">{value}</span>
</div>

<style>
  .metric {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  .metric-label {
    font-size: 0.85rem;
    color: color-mix(in srgb, var(--text) 60%, transparent);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .metric-value {
    font-size: 2.827rem; /* Equivalent to h3 in scale */
    color: var(--text);
    line-height: 1;
  }
</style>
