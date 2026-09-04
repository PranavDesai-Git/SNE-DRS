<script>
  import Button from './Button.svelte';
  import TextInput from './TextInput.svelte';
  import Icon from './Icon.svelte';

  export let label = "Location";
  export let lat = "";
  export let lng = "";
  
  let isLoading = false;

  function getLocation() {
    isLoading = true;
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          lat = position.coords.latitude.toFixed(6);
          lng = position.coords.longitude.toFixed(6);
          isLoading = false;
        },
        (error) => {
          console.error("Error getting location", error);
          isLoading = false;
          alert("Unable to retrieve location.");
        }
      );
    } else {
      isLoading = false;
      alert("Geolocation is not supported by this browser.");
    }
  }
</script>

<div class="location-picker">
  {#if label}
    <label>{label}</label>
  {/if}
  <div class="controls">
    <TextInput placeholder="Latitude" bind:value={lat} />
    <TextInput placeholder="Longitude" bind:value={lng} />
    <Button on:click={getLocation} disabled={isLoading}>
      {#if isLoading}
        Getting...
      {:else}
        <Icon name="search" size="18" /> GPS
      {/if}
    </Button>
  </div>
</div>

<style>
  .location-picker {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  label {
    font-size: 0.85rem;
    color: color-mix(in srgb, var(--text) 60%, transparent);
  }
  .controls {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }
</style>
