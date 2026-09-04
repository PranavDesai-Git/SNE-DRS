<script>
  export let src = "";
  export let alt = "Image";
  export let caption = "";
  
  let loaded = false;
  let error = false;
</script>

<div class="gov-image-block">
  {#if !loaded && !error}
    <div class="image-placeholder">Loading...</div>
  {/if}
  
  {#if error}
    <div class="image-placeholder error">Failed to load image</div>
  {:else}
    <img 
      {src} 
      {alt} 
      on:load={() => loaded = true}
      on:error={() => error = true}
      class:visible={loaded}
    />
  {/if}
  
  {#if caption}
    <div class="caption">{caption}</div>
  {/if}
</div>

<style>
  .gov-image-block {
    display: flex;
    flex-direction: column;
    border: 1px solid var(--secondary);
    border-radius: 4px;
    overflow: hidden;
    background-color: color-mix(in srgb, var(--background), black 2%);
  }

  img {
    width: 100%;
    display: block;
    object-fit: cover;
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  img.visible {
    opacity: 1;
  }

  .image-placeholder {
    padding: 3rem 1rem;
    text-align: center;
    font-size: 0.85rem;
    color: color-mix(in srgb, var(--text) 60%, transparent);
  }
  
  .caption {
    padding: 0.5rem;
    font-size: 0.8rem;
    color: var(--text);
    background-color: var(--background);
    border-top: 1px solid var(--secondary);
    text-align: center;
  }
</style>
