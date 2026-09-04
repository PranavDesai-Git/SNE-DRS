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
  
  $: templateText = `${habitationName} is flagged ${tier} priority, hazard score ${score}/1.0. Nearest viable relocation site: ${siteName} (capacity ${capacity}, ${distance}km away).`;
  
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
  
  <slot></slot>
</Card>

<style>
  .summary {
    margin-top: 1rem;
    margin-bottom: 1rem;
  }
</style>
