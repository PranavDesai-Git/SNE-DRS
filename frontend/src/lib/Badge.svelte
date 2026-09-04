<script>
  // tier can be: 'immediate', 'moderate', 'safe', 'habitation'
  export let tier = 'safe'; 
  export let isBand = false;
</script>

<span class="badge {tier}" class:is-band={isBand}>
  <slot />
</span>

<style>
  .badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    border-radius: 4px; /* Matches the button rounding */
    font-family: 'Poppins', sans-serif;
    font-size: 0.707rem; /* Ties into your typography scale (small) */
    text-transform: uppercase;
    color: white; /* Ensures high contrast against dark risk colors */
    letter-spacing: 0.05em;
  }

  .is-band {
    display: block;
    width: 100%;
    border-radius: 0;
    text-align: center;
    padding: 0.4rem;
  }

  /* Map specifically to your provided risk colors */
  .immediate { background-color: var(--color-redzone); }
  .moderate { background-color: var(--color-moderate); }
  .safe { background-color: var(--color-safe); }
  .habitation { background-color: var(--color-habitation); }
</style>
