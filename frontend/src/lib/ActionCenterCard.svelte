<script>
  import Card from './Card.svelte';
  import Badge from './Badge.svelte';
  import Callout from './Callout.svelte';
  
  export let habitationName = "Unknown";
  export let tier = "moderate";
  export let score = "0.0";
  export let siteName = "None";
  export let capacity = 0;
  export let distance = 0;
  
  export let slopeScore = 0;
  export let rainfallScore = 0;
  export let twiScore = 0;
  export let landcoverScore = 0;
  
  export let pctElderly = 0;
  export let pctChildren = 0;
  export let population = 0;
  
  export let currentRations = 0;
  export let cots = 0;
  export let medicalKits = 0;
  
  $: templateText = `${habitationName} is flagged ${tier} priority, hazard score ${score}/10. Nearest viable relocation site: ${siteName} (capacity ${capacity}, ${distance}km away).`;
  
  // Decide callout type based on tier
  $: calloutType = tier === 'immediate' ? 'alert' : 'info';
</script>

<Card title={`Action Center: ${habitationName}`}>
  <Badge slot="band" {tier} isBand>{tier.toUpperCase()} PRIORITY</Badge>
  
  <div class="summary">
    <Callout type={calloutType}>
      {templateText}
    </Callout>
  </div>
  
  <div class="breakdown">
    <small>Hazard & Demographics:</small>
    <div class="metrics">
      <div class="metric">Slope: {slopeScore}</div>
      <div class="metric">Rainfall: {rainfallScore}</div>
      <div class="metric">Elderly: {pctElderly}%</div>
      <div class="metric">Children: {pctChildren}%</div>
    </div>
  </div>

  <div class="breakdown">
    <div style="display: flex; justify-content: space-between;">
      <small>Site Logistics (vs {population} Pop):</small>
      <button style="border: none; background: none; cursor: pointer; color: var(--primary-dark); text-decoration: underline;" on:click={() => window.print()}>Print Brief</button>
    </div>
    <div class="metrics">
      <div class="metric" class:deficit={currentRations < population}>Rations: {currentRations}</div>
      <div class="metric" class:deficit={cots < population}>Cots: {cots}</div>
      <div class="metric">Med Kits: {medicalKits}</div>
    </div>
  </div>
  
  <slot></slot>
</Card>

<style>
  .deficit {
    color: #87231E; /* Redzone color for deficits, no bolding */
  }
  .summary {
    margin-top: 1rem;
    margin-bottom: 1rem;
  }
  .breakdown {
    margin-top: 1rem;
    padding: 0.5rem;
    background: rgba(255, 255, 255, 0.4);
    border-radius: 4px;
    margin-bottom: 1rem;
  }
  .metrics {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.5rem;
    font-size: 0.8rem;
    margin-top: 0.5rem;
  }
</style>
