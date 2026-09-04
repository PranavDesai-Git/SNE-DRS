<script>
  export let type = "info"; // info, alert
</script>

<div class="gov-callout {type}">
  <slot />
</div>

<style>
  .gov-callout {
    padding: 1rem;
    border-radius: 4px;
    font-size: 0.95rem;
    line-height: 1.5;
  }
  
  /* According to the DESIGN_SYSTEM.md, we must not mix pastels or use external shadows/borders on cards. 
     This is an inline callout, so we keep it simple. */
  .info {
    background-color: color-mix(in srgb, var(--primary), transparent 85%);
    color: var(--text);
    border: 1px solid var(--primary);
  }

  .alert {
    background-color: color-mix(in srgb, var(--color-moderate), transparent 85%);
    border: 1px solid var(--color-moderate);
  }
</style>
