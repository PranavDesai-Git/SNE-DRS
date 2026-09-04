<script>
  import Icon from './Icon.svelte';
  export let title = "No Data Found";
  export let description = "";
  export let icon = "search";
</script>

<div class="gov-empty-state">
  <div class="icon-wrapper">
    <Icon name={icon} size="32" />
  </div>
  <h4>{title}</h4>
  {#if description}
    <p>{description}</p>
  {/if}
  <slot />
</div>

<style>
  .gov-empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 3rem 1rem;
    text-align: center;
    border: 1px dashed var(--secondary);
    border-radius: 4px;
    background-color: color-mix(in srgb, var(--background), black 2%);
  }

  .icon-wrapper {
    color: var(--secondary);
    margin-bottom: 1rem;
  }

  h4 {
    margin: 0 0 0.5rem 0;
    color: var(--text);
  }

  p {
    font-size: 0.9rem;
    color: color-mix(in srgb, var(--text) 60%, transparent);
    margin-bottom: 1rem;
    max-width: 300px;
  }
</style>
