<script>
  export let paths = []; // Array of { name: string, href: string }
</script>

<nav class="gov-breadcrumbs">
  {#each paths as path, i}
    <a href={path.href || '#'} class="breadcrumb-link" class:active={i === paths.length - 1}>
      {path.name}
    </a>
    {#if i < paths.length - 1}
      <span class="separator">/</span>
    {/if}
  {/each}
</nav>

<style>
  .gov-breadcrumbs {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.85rem;
  }

  .breadcrumb-link {
    color: var(--text);
    text-decoration: none;
    transition: opacity 0.2s;
  }

  .breadcrumb-link:not(.active):hover {
    text-decoration: underline;
  }

  .breadcrumb-link.active {
    font-weight: 400; /* Strict adherence to typography rules */
    opacity: 0.8;
    pointer-events: none;
  }

  .separator {
    color: color-mix(in srgb, var(--text) 60%, transparent);
  }
</style>
