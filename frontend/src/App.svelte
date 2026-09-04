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

  function handleSelectHabitation(hab) {
    selectedHabitation = hab;
    if (sites.length > 0) {
      assignedSite = sites.find(s => s.capacity >= hab.population) || sites[0] || {};
      console.log(`[Frontend] Re-assigned ${selectedHabitation.name} to ${assignedSite.name}`);
    }
  }

  function exportCSV() {
    const header = "id,name,district,lat,lng,population,risk_score,tier,slope_score,twi_score,landcover_score,rainfall_score\n";
    const rows = habitations.map(h => 
      `${h.id},${h.name},${h.district},${h.lat},${h.lng},${h.population},${h.risk_score},${h.tier},${h.slope_score},${h.twi_score},${h.landcover_score},${h.rainfall_score}`
    ).join("\n");
    const blob = new Blob([header + rows], { type: 'text/csv' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'habitations_priority.csv';
    a.click();
    URL.revokeObjectURL(url);
  }

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

  let reports = [];
  let currentRoute = null;

  async function loadReports() {
    try {
      const res = await fetch('http://localhost:8080/reports');
      if (res.ok) reports = await res.json() || [];
    } catch (e) { console.error(e); }
  }

  async function loadRoute(habId) {
    try {
      const res = await fetch(`http://localhost:8080/routes?habitation_id=${habId}`);
      if (res.ok) currentRoute = await res.json() || null;
    } catch (e) { console.error(e); }
  }

  $: if (selectedMode === 'm2' && selectedHabitation) {
    loadRoute(selectedHabitation.id);
  }
  
  $: if (selectedMode === 'm3') {
    loadReports();
  }

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
          <div style="display: flex; gap: 0.5rem;">
            <Button on:click={exportCSV}>Export CSV</Button>
            <Button on:click={() => showDataProvenance = true}>Data Provenance</Button>
          </div>
        </div>
        
        <DataTable 
          headers={tableHeaders} 
          data={habitations} 
          sortable={true} 
          onRowClick={handleSelectHabitation}
          selectedId={selectedHabitation?.id}
        />

        {#if selectedHabitation && assignedSite}
        <div style="margin-top: 1rem;">
          <ActionCenterCard 
            habitationName={selectedHabitation.name}
            tier={selectedHabitation.tier.toLowerCase()}
            score={selectedHabitation.risk_score.toString()}
            siteName={assignedSite.name || 'None'}
            capacity={assignedSite.capacity || 0}
            distance={assignedSite.distance_to_road_km || 0}
            slopeScore={selectedHabitation.slope_score}
            rainfallScore={selectedHabitation.rainfall_score}
            twiScore={selectedHabitation.twi_score}
            landcoverScore={selectedHabitation.landcover_score}
            pctElderly={selectedHabitation.pct_elderly}
            pctChildren={selectedHabitation.pct_children}
            population={selectedHabitation.population}
            currentRations={assignedSite.current_rations}
            cots={assignedSite.cots}
            medicalKits={assignedSite.medical_kits}
          />
        </div>
        {/if}
      {:else if selectedMode === 'm2'}
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem;">
          <h4>Pre-Disaster Routing</h4>
        </div>
        <p class="subtitle">Select a habitation to view evacuation routes.</p>
        
        <DataTable 
          headers={tableHeaders} 
          data={habitations} 
          sortable={true} 
          onRowClick={handleSelectHabitation}
          selectedId={selectedHabitation?.id}
        />

        {#if currentRoute && selectedHabitation && assignedSite}
        <div style="margin-top: 1rem; padding: 1rem; background: rgba(255, 255, 255, 0.4); border-radius: 4px;">
          <h5>Evacuation Route Generated</h5>
          <p><strong>From:</strong> {selectedHabitation.name}</p>
          <p><strong>To:</strong> {assignedSite.name}</p>
          <p><strong>Status:</strong> Risk zones avoided</p>
          <Button fullWidth style="margin-top: 0.5rem;">Deploy Resources</Button>
        </div>
        {/if}

      {:else if selectedMode === 'm3'}
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem;">
          <h4>Active Incidents</h4>
          <Button on:click={loadReports}>Refresh</Button>
        </div>
        <p class="subtitle">Real-time citizen reports and SOS signals.</p>
        
        <div style="display: flex; flex-direction: column; gap: 0.5rem; margin-top: 1rem;">
          {#each reports as report}
            <div style="padding: 1rem; border: 1px solid var(--color-redzone); border-radius: 4px; background: rgba(135, 35, 30, 0.1);">
              <div style="display: flex; justify-content: space-between;">
                <strong>{report.category.toUpperCase()}</strong>
                <small>{new Date(report.reported_at).toLocaleString()}</small>
              </div>
              <p style="margin-top: 0.5rem;">{report.description}</p>
              <div style="margin-top: 0.5rem; display: flex; justify-content: space-between; font-size: 0.85rem;">
                <span>Status: {report.status}</span>
                <Button>Dispatch Team</Button>
              </div>
            </div>
          {/each}
          {#if reports.length === 0}
            <p>No active reports.</p>
          {/if}
        </div>
      {/if}
    </div>
  </aside>

  <!-- Resizer Hotdog -->
  <div class="resizer" style="left: {sidebarWidth}px;" on:mousedown={startResize} role="separator" tabindex="0">
    <div class="hotdog"></div>
  </div>
  {/if}

  <!-- Right Map Pane -->
  <section class="map-container">
    <button class="hamburger" on:click={() => sidebarCollapsed = !sidebarCollapsed}>
      ☰
    </button>
    <Map 
      {riskZones} 
      {sites} 
      {habitations}
      {reports}
      {currentRoute}
      {selectedMode}
      onHabitationClick={handleSelectHabitation} 
    />
    <div class="legend-wrapper">
      <MapLegend />
    </div>
  </section>

  <!-- Drawer placed outside so it is not trapped by backdrop-filter -->
  <Drawer bind:isOpen={showDataProvenance} title="Data Provenance" position="right">
    <p>Elevation Source: CartoDEM via Bhoonidhi</p>
    <p>Landslide Inventory: GSI NLSM</p>
    <p>Flood Hazard: Bhuvan Web Services</p>
    <p>Population: Census 2011 (Growth-adjusted)</p>
    <div style="margin-top: 2rem;">
      <Button fullWidth on:click={() => showDataProvenance = false}>Close</Button>
    </div>
  </Drawer>
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
    position: absolute;
    top: 0;
    height: 100vh;
    width: 12px;
    background-color: color-mix(in srgb, var(--background) 55%, transparent);
    backdrop-filter: blur(24px);
    -webkit-backdrop-filter: blur(24px);
    cursor: col-resize;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1001;
    border-right: 1px solid rgba(255, 255, 255, 0.5);
  }

  .resizer:hover, .resizer:active {
    background-color: color-mix(in srgb, var(--background) 45%, transparent);
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
    left: calc(var(--sidebar-width, 420px) + 2rem);
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

  @media print {
    :global(body) {
      background: white !important;
      color: black !important;
    }
    .resizer, .hamburger, .mode-selector, .legend-wrapper {
      display: none !important;
    }
    .app-container {
      display: block !important;
      height: auto !important;
    }
    .sidebar {
      width: 100% !important;
      border: none !important;
      height: auto !important;
      overflow: visible !important;
    }
    .map-container {
      width: 100% !important;
      height: 400px !important;
      position: relative !important;
      border: 1px solid #ccc;
      margin-top: 2rem;
    }
    :global(.leaflet-control-container) {
      display: none !important;
    }
  }
</style>
