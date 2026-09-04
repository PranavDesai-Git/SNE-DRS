<script>
  import Button from './Button.svelte';
  
  export let label = "";
  export let accept = "image/*";
  export let files = null;
  
  let fileInput;

  function triggerSelect() {
    fileInput.click();
  }

  function handleFileChange(event) {
    files = event.target.files;
  }
</script>

<div class="input-group">
  {#if label}
    <label>{label}</label>
  {/if}
  
  <div class="file-area">
    <Button type="button" on:click={triggerSelect}>
      Choose File
    </Button>
    <span class="file-name">
      {#if files && files.length > 0}
        {files[0].name}
      {:else}
        No file chosen
      {/if}
    </span>
  </div>

  <input 
    type="file" 
    bind:this={fileInput}
    on:change={handleFileChange}
    {accept}
    class="hidden-input" 
  />
</div>

<style>
  .input-group {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }
  
  label {
    font-size: 0.85rem;
    font-weight: 400;
  }

  .file-area {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .file-name {
    font-size: 0.9rem;
    color: color-mix(in srgb, var(--text) 60%, transparent);
  }

  .hidden-input {
    display: none;
  }
</style>
