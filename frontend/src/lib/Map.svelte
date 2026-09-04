<script>
  import { onMount } from 'svelte';
  
  export let lat = 25.5788; // Default NE India roughly (Meghalaya)
  export let lng = 91.8933; // Shillong
  export let zoom = 8;
  
  let mapContainer;
  let map;

  onMount(async () => {
    // Dynamic import to avoid SSR issues if this were SvelteKit (though this is Vite SPA, good practice)
    const L = await import('leaflet');
    await import('leaflet/dist/leaflet.css');

    map = L.map(mapContainer).setView([lat, lng], zoom);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);

    return () => {
      map.remove();
    };
  });
</script>

<div class="map-wrapper" bind:this={mapContainer}></div>

<style>
  .map-wrapper {
    width: 100%;
    height: 100%;
    min-height: 400px;
    background-color: var(--secondary);
    border: 1px solid var(--secondary);
    border-radius: 4px;
    z-index: 1;
  }
</style>
