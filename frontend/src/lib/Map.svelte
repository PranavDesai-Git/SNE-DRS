<script>
  import { onMount } from 'svelte';
  
  export let lat = 26.4521; // Barpeta roughly
  export let lng = 90.9871;
  export let zoom = 8;
  export let riskZones = null;
  export let sites = [];
  export let habitations = [];
  export let reports = [];
  export let currentRoute = null;
  export let selectedMode = 'm1';
  export let onHabitationClick = null;
  
  let mapContainer;
  let map;
  let leafletBase;
  let geoJsonLayer;
  let pulseLayer;
  let waypointMarker;
  let siteMarkers = [];
  let habitationMarkers = [];
  let reportMarkers = [];
  let routeLayer = null;

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
    renderHabitations();

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

  $: if (habitations) {
    renderHabitations();
  }

  $: if (currentRoute || selectedMode) {
    renderRoute();
  }

  $: if (reports || selectedMode) {
    renderReports();
  }

  function renderRoute() {
    if (!map || !leafletBase) return;
    if (routeLayer) {
      routeLayer.remove();
      routeLayer = null;
    }
    if (selectedMode !== 'm2' || !currentRoute) return;

    routeLayer = leafletBase.geoJSON(currentRoute, {
      style: {
        color: '#156932',
        weight: 4,
        dashArray: '10, 10',
        lineCap: 'round'
      }
    }).addTo(map);
  }

  function renderReports() {
    if (!map || !leafletBase) return;
    reportMarkers.forEach(m => m.remove());
    reportMarkers = [];

    if (selectedMode !== 'm3' || !reports) return;

    reports.forEach(rep => {
      const icon = leafletBase.divIcon({
        className: 'custom-report-icon',
        html: `<div style="background-color: #87231E; color: white; border-radius: 50%; width: 20px; height: 20px; display: flex; align-items: center; justify-content: center; font-size: 14px; border: 2px solid white; box-shadow: 0 2px 4px rgba(0,0,0,0.4);">!</div>`,
        iconSize: [20, 20],
        iconAnchor: [10, 10]
      });

      const m = leafletBase.marker([rep.lat, rep.lng], { icon: icon })
        .bindPopup(`<h5>${rep.category.toUpperCase()}</h5><p>${rep.description}</p><em>Status: ${rep.status}</em>`)
        .addTo(map);

      reportMarkers.push(m);
    });
  }

  function renderHabitations() {
    if (!map || !leafletBase || !habitations) return;
    
    habitationMarkers.forEach(m => m.remove());
    habitationMarkers = [];

    habitations.forEach(hab => {
      const habIcon = leafletBase.divIcon({
        className: 'custom-hab-icon',
        html: `<div style="background-color: #0f1510; border: 2px solid white; border-radius: 50%; width: 14px; height: 14px; box-shadow: 0 1px 3px rgba(0,0,0,0.4);"></div>`,
        iconSize: [14, 14],
        iconAnchor: [7, 7]
      });

      const m = leafletBase.marker([hab.lat, hab.lng], { icon: habIcon })
        .bindTooltip(`<h5>${hab.name}</h5>Pop: ${hab.population}`)
        .addTo(map)
        .on('click', () => {
          if (onHabitationClick) {
            onHabitationClick(hab);
          }
        });

      habitationMarkers.push(m);
    });
  }

  function renderSites() {
    if (!map || !leafletBase || !sites) return;
    
    // Clear old site markers
    siteMarkers.forEach(m => m.remove());
    siteMarkers = [];

    sites.forEach(site => {
      const safeIcon = leafletBase.divIcon({
        className: 'custom-safe-icon',
        html: '<div style="background-color: #156932; color: white; border-radius: 50%; width: 24px; height: 24px; display: flex; align-items: center; justify-content: center; font-size: 18px; border: 2px solid white; box-shadow: 0 2px 4px rgba(0,0,0,0.3);">+</div>',
        iconSize: [24, 24],
        iconAnchor: [12, 12]
      });

      const m = leafletBase.marker([site.lat, site.lng], { icon: safeIcon })
        .bindPopup(`<h5>Safe Relocation Site</h5>${site.name}<br/>Capacity: ${site.capacity}`)
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
          const hazardType = feature.properties.hazard_type || 'Unknown';
          const icons = {
            'landslide': '⛰️ Landslide',
            'flood': '🌊 Flood',
            'earthquake': '🌍 Earthquake',
            'tornado': '🌪️ Tornado'
          };
          const emojiMap = {
            'landslide': '⛰️',
            'flood': '🌊',
            'earthquake': '🌍',
            'tornado': '🌪️'
          };
          const hazardDisplay = icons[hazardType.toLowerCase()] || hazardType;
          const emoji = emojiMap[hazardType.toLowerCase()] || '⚠️';

          layer.bindPopup(
            `<h5>Zone ID: ${feature.properties.zone_id}</h5>` +
            `<strong>Hazard:</strong> ${hazardDisplay}<br/>` +
            `Tier: ${feature.properties.tier}<br/>` +
            `Risk Score: ${feature.properties.risk_score}`
          );

          layer.bindTooltip(emoji, {
            permanent: true,
            direction: 'center',
            className: 'hazard-emoji-tooltip'
          });
        }
      }
    }).addTo(map);

    if (pulseLayer) {
      pulseLayer.remove();
    }
    
    pulseLayer = leafletBase.geoJSON(riskZones, {
      style: function(feature) {
        let speedClass = 'pulse-speed-moderate';
        if (feature.properties && feature.properties.tier) {
          switch(feature.properties.tier.toLowerCase()) {
            case 'immediate': speedClass = 'pulse-speed-immediate'; break;
            case 'short-term': speedClass = 'pulse-speed-short-term'; break;
            case 'medium-term': speedClass = 'pulse-speed-medium-term'; break;
            case 'moderate': speedClass = 'pulse-speed-moderate'; break;
            case 'safe': speedClass = 'pulse-speed-safe'; break;
          }
        }
        
        return {
          color: getColor(feature.properties.tier),
          weight: 2,
          fillColor: getColor(feature.properties.tier),
          className: `pulse-ring-path ${speedClass}`
        };
      },
      interactive: false
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
        .bindPopup("<h5>High Risk Cluster</h5>Click to zoom in")
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
  
  :global(.hazard-emoji-tooltip) {
    background: transparent !important;
    border: none !important;
    box-shadow: none !important;
    font-size: 24px;
    text-shadow: 0px 0px 4px rgba(255, 255, 255, 0.8);
  }
</style>
