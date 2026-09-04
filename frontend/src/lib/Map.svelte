<script>
  import { onMount } from 'svelte';
  
  export let lat = 26.4521; // Barpeta roughly
  export let lng = 90.9871;
  export let zoom = 8;
  export let riskZones = null;
  export let sites = [];
  
  let mapContainer;
  let map;
  let leafletBase;
  let geoJsonLayer;
  let waypointMarker;
  let siteMarkers = [];

  function getColor(tier) {
    if (!tier) return '#8F7518';
    switch(tier.toLowerCase()) {
      case 'immediate': return '#87231E';
      case 'short-term': return '#8F7518';
      case 'medium-term': return '#8F7518';
      case 'moderate': return '#8F7518';
      case 'safe': return '#354F3E';
      default: return '#8F7518';
    }
  }

  onMount(async () => {
    leafletBase = await import('leaflet');
    await import('leaflet/dist/leaflet.css');

    map = leafletBase.map(mapContainer).setView([lat, lng], zoom);

    leafletBase.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; OpenStreetMap contributors'
    }).addTo(map);

    map.on('zoomend', updateWaypoint);

    renderPolygons();
    renderSites();

    return () => {
      if (map) map.remove();
    };
  });

  $: if (riskZones) {
    renderPolygons();
  }

  $: if (sites) {
    renderSites();
  }

  function renderSites() {
    if (!map || !leafletBase || !sites) return;
    
    // Clear old site markers
    siteMarkers.forEach(m => m.remove());
    siteMarkers = [];

    sites.forEach(site => {
      const safeIcon = leafletBase.divIcon({
        className: 'custom-safe-icon',
        html: '<div style="background-color: #156932; color: white; border-radius: 50%; width: 24px; height: 24px; display: flex; align-items: center; justify-content: center; font-weight: bold; font-size: 18px; border: 2px solid white; box-shadow: 0 2px 4px rgba(0,0,0,0.3);">+</div>',
        iconSize: [24, 24],
        iconAnchor: [12, 12]
      });

      const m = leafletBase.marker([site.lat, site.lng], { icon: safeIcon })
        .bindPopup(`<strong>Safe Relocation Site</strong><br/>${site.name}<br/>Capacity: ${site.capacity}`)
        .addTo(map);

      siteMarkers.push(m);
    });
  }

  function renderPolygons() {
    if (!map || !leafletBase || !riskZones) return;

    if (geoJsonLayer) {
      geoJsonLayer.remove();
    }

    geoJsonLayer = leafletBase.geoJSON(riskZones, {
      style: function(feature) {
        return {
          color: getColor(feature.properties.tier),
          weight: 2,
          opacity: 1,
          fillColor: getColor(feature.properties.tier),
          fillOpacity: 0.4
        };
      },
      onEachFeature: function(feature, layer) {
        if (feature.properties) {
          layer.bindPopup(
            `<strong>Zone ID:</strong> ${feature.properties.zone_id}<br/>` +
            `<strong>Tier:</strong> ${feature.properties.tier}<br/>` +
            `<strong>Risk Score:</strong> ${feature.properties.risk_score}`
          );
        }
      }
    }).addTo(map);

    if (riskZones.features && riskZones.features.length > 0) {
      map.fitBounds(geoJsonLayer.getBounds());
      updateWaypoint();
    }
  }

  function updateWaypoint() {
    if (!map || !leafletBase || !geoJsonLayer) return;

    const currentZoom = map.getZoom();
    const thresholdZoom = 10; // Switch to waypoint when zoomed out beyond 10

    if (currentZoom < thresholdZoom) {
      if (!waypointMarker) {
        // Place the waypoint at the center of the hazard zones
        const center = geoJsonLayer.getBounds().getCenter();
        
        waypointMarker = leafletBase.circleMarker(center, {
          radius: 8,
          color: '#0f1510',
          weight: 2,
          fillColor: '#87231E', // Red zone color
          fillOpacity: 1
        })
        .bindPopup("<strong>High Risk Cluster</strong><br/>Click to zoom in")
        .addTo(map)
        .on('click', () => {
          map.setView(center, 12); // Zoom in past threshold
        });
      }
    } else {
      if (waypointMarker) {
        waypointMarker.remove();
        waypointMarker = null;
      }
    }
  }
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
