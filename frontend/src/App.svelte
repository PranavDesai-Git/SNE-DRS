<script>
  import { onMount } from 'svelte';
  import Map from './lib/Map.svelte';
  import ButtonGroup from './lib/ButtonGroup.svelte';
  import MapLegend from './lib/MapLegend.svelte';
  import DataTable from './lib/DataTable.svelte';
  import ActionCenterCard from './lib/ActionCenterCard.svelte';
  import Button from './lib/Button.svelte';
  import Drawer from './lib/Drawer.svelte';

  let modeOptions = [
    { label: 'Mode 1: Baseline', value: 'm1' },
    { label: 'Mode 2: Pre-Disaster', value: 'm2' },
    { label: 'Mode 3: Active', value: 'm3' }
  ];
  let selectedMode = 'm1';

  let habitations = [];
  let sites = [];
  let riskZones = null;

  let selectedHabitation = null;
  let assignedSite = null;

  let sidebarWidth = 420;
  let sidebarCollapsed = false;
  let isResizing = false;

  function startResize() { isResizing = true; }
  function stopResize() { isResizing = false; }
  function onMouseMove(e) {
    if (isResizing) {
      sidebarWidth = Math.max(300, Math.min(e.clientX, 800));
    }
  }

  onMount(async () => {
    try {
      const [habRes, sitesRes, zonesRes] = await Promise.all([
        fetch('http://localhost:8080/habitations'),
        fetch('http://localhost:8080/sites'),
        fetch('http://localhost:8080/risk-zones')
      ]);

      if (habRes.ok) habitations = await habRes.json() || [];
      console.log(`[Frontend] Fetched ${habitations.length} habitations`);

      if (sitesRes.ok) sites = await sitesRes.json() || [];
      console.log(`[Frontend] Fetched ${sites.length} sites`);

      if (zonesRes.ok) riskZones = await zonesRes.json() || null;
      console.log(`[Frontend] Fetched risk zones:`, riskZones ? riskZones.features?.length + " features" : "null");

      if (habitations.length > 0) {
        // Default to the highest risk habitation
        selectedHabitation = habitations[0];
        // Match with a site that has enough capacity
        assignedSite = sites.find(s => s.capacity >= selectedHabitation.population) || sites[0] || {};
        console.log(`[Frontend] Assigned ${selectedHabitation.name} to ${assignedSite.name}`);
      }

    } catch (e) {
      console.error("Failed to fetch data from localhost:8080", e);
    }
  });

  let tableHeaders = [
    { key: 'name', label: 'Habitation', sortable: true },
    { key: 'district', label: 'District', sortable: true },
    { key: 'tier', label: 'Priority', sortable: true },
    { key: 'risk_score', label: 'Score', sortable: true }
  ];

  let showDataProvenance = false;
</script>

<svelte:window on:mousemove={onMouseMove} on:mouseup={stopResize} />

<main class="app-container" style="cursor: {isResizing ? 'col-resize' : 'default'};">
  <!-- Left Sidebar -->
  {#if !sidebarCollapsed}
  <aside class="sidebar" style="width: {sidebarWidth}px; flex-shrink: 0;">
    <div class="sidebar-header">
      <h3>GIS Based Disaster Prevention System</h3>
      <p class="subtitle">Disaster Relocation Planning</p>
    </div>
    
    <div class="mode-selector">
      <ButtonGroup options={modeOptions} bind:selected={selectedMode} />
    </div>

    <div class="sidebar-content">
      {#if selectedMode === 'm1'}
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem;">
          <h4>Habitation Risk Ranking</h4>
          <Button on:click={() => showDataProvenance = true}>Data Provenance</Button>
        </div>
        
        <DataTable headers={tableHeaders} data={habitations} sortable={true} />

        {#if selectedHabitation && assignedSite}
        <div style="margin-top: 1rem;">
          <ActionCenterCard 
            habitationName={selectedHabitation.name}
            tier={selectedHabitation.tier.toLowerCase()}
            score={selectedHabitation.risk_score.toString()}
            siteName={assignedSite.name || 'None'}
            capacity={assignedSite.capacity || 0}
            distance={assignedSite.distance_to_road_km || 0}
          />
        </div>
        {/if}

        <Drawer bind:isOpen={showDataProvenance} title="Data Provenance" position="right">
          <p><strong>Elevation Source:</strong> CartoDEM via Bhoonidhi</p>
          <p><strong>Landslide Inventory:</strong> GSI NLSM</p>
          <p><strong>Flood Hazard:</strong> Bhuvan Web Services</p>
          <p><strong>Population:</strong> Census 2011 (Growth-adjusted)</p>
          <div style="margin-top: 2rem;">
            <Button fullWidth on:click={() => showDataProvenance = false}>Close</Button>
          </div>
        </Drawer>
      {:else if selectedMode === 'm2'}
        <h4>Pre-Disaster Routing</h4>
        <p class="subtitle">Waiting for threshold trigger...</p>
        <!-- Mode 2 Routing & Capacity matching will go here -->
      {:else if selectedMode === 'm3'}
        <h4>Active Incidents</h4>
        <p class="subtitle">Awaiting citizen reports...</p>
        <!-- Mode 3 Real-time citizen reports will go here -->
      {/if}
    </div>
  </aside>

  <!-- Resizer Hotdog -->
  <div class="resizer" on:mousedown={startResize} role="separator" tabindex="0">
    <div class="hotdog"></div>
  </div>
  {/if}

  <!-- Right Map Pane -->
  <section class="map-container">
    <button class="hamburger" on:click={() => sidebarCollapsed = !sidebarCollapsed}>
      ☰
    </button>
    <Map {riskZones} {sites} />
    <div class="legend-wrapper">
      <MapLegend />
    </div>
  </section>
</main>

<style>
  .sidebar-header {
    padding: 1.5rem;
    border-bottom: 1px solid var(--secondary);
  }

  .subtitle {
    font-size: 0.85rem;
    color: color-mix(in srgb, var(--text) 60%, transparent);
  }

  .mode-selector {
    padding: 1.5rem;
    border-bottom: 1px solid var(--secondary);
  }

  .sidebar-content {
    flex: 1;
    overflow-y: auto;
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .legend-wrapper {
    position: absolute;
    bottom: 2rem;
    right: 2rem;
    z-index: 1000;
  }

  .resizer {
    width: 12px;
    background-color: color-mix(in srgb, var(--background), var(--secondary) 20%);
    cursor: col-resize;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1001;
    border-right: 1px solid var(--secondary);
  }

  .resizer:hover, .resizer:active {
    background-color: color-mix(in srgb, var(--background), var(--secondary) 40%);
  }

  .hotdog {
    width: 4px;
    height: 30px;
    background-color: var(--primary-dark);
    border-radius: 2px;
  }

  .hamburger {
    position: absolute;
    top: 1rem;
    left: 1rem;
    z-index: 2000;
    background: var(--background);
    border: 1px solid var(--secondary);
    border-radius: 4px;
    padding: 0.5rem 0.75rem;
    cursor: pointer;
    font-size: 1.2rem;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
    color: var(--text);
  }

  .hamburger:hover {
    background: color-mix(in srgb, var(--background), var(--primary) 20%);
  }
</style>
