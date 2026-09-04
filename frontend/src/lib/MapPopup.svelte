<script>
  import Badge from './Badge.svelte';
  import Divider from './Divider.svelte';
  import Button from './Button.svelte';

  export let title = "Unknown Location";
  export let tier = "moderate"; // immediate, moderate, safe, habitation
  export let population = 0;
  export let hazardScore = 0.0;
  export let onAction = () => {};
</script>

<div class="map-popup">
  <div class="header">
    <strong>{title}</strong>
    <Badge {tier}>{tier.charAt(0).toUpperCase() + tier.slice(1)} Risk</Badge>
  </div>
  
  <Divider />
  
  <div class="body">
    <p>Population: <strong>{population}</strong></p>
    <p>Hazard Score: <strong>{hazardScore}/1.0</strong></p>
  </div>

  <div class="actions">
    <Button fullWidth on:click={onAction}>View Details</Button>
  </div>
</div>

<style>
  .map-popup {
    font-family: 'Poppins', sans-serif;
    color: var(--text);
    min-width: 220px;
  }
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    margin-bottom: 0.5rem;
  }
  .body {
    margin: 0.5rem 0;
    font-size: 0.85rem;
  }
  .body p {
    margin: 0.2rem 0;
  }
  .actions {
    margin-top: 1rem;
  }
</style>
