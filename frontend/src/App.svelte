<script>
  import Button from './lib/Button.svelte';
  import Badge from './lib/Badge.svelte';
  import Card from './lib/Card.svelte';
  import TextInput from './lib/TextInput.svelte';
  import Select from './lib/Select.svelte';
  import Slider from './lib/Slider.svelte';
  import Switch from './lib/Switch.svelte';
  import Checkbox from './lib/Checkbox.svelte';
  import Spinner from './lib/Spinner.svelte';
  import Tabs from './lib/Tabs.svelte';
  import Modal from './lib/Modal.svelte';
  import Notification from './lib/Notification.svelte';
  import ProgressBar from './lib/ProgressBar.svelte';
  import Tooltip from './lib/Tooltip.svelte';
  import Accordion from './lib/Accordion.svelte';
  import Divider from './lib/Divider.svelte';
  import Textarea from './lib/Textarea.svelte';
  import SearchInput from './lib/SearchInput.svelte';
  import FileInput from './lib/FileInput.svelte';
  import Metric from './lib/Metric.svelte';
  import Tag from './lib/Tag.svelte';
  import Callout from './lib/Callout.svelte';
  import Skeleton from './lib/Skeleton.svelte';
  import DataTable from './lib/DataTable.svelte';
  import Drawer from './lib/Drawer.svelte';
  import Pagination from './lib/Pagination.svelte';
  import MapLegend from './lib/MapLegend.svelte';
  import RadioGroup from './lib/RadioGroup.svelte';
  import DropdownMenu from './lib/DropdownMenu.svelte';
  import Breadcrumbs from './lib/Breadcrumbs.svelte';
  import Icon from './lib/Icon.svelte';
  import EmptyState from './lib/EmptyState.svelte';
  import ButtonGroup from './lib/ButtonGroup.svelte';
  import Popover from './lib/Popover.svelte';
  import ImageBlock from './lib/ImageBlock.svelte';
  import Stepper from './lib/Stepper.svelte';
  import Map from './lib/Map.svelte';
  import MapPopup from './lib/MapPopup.svelte';
  import LocationPicker from './lib/LocationPicker.svelte';
  import ActionCenterCard from './lib/ActionCenterCard.svelte';

  let showModal = false;
  let showNotification = false;
  let showDrawer = false;
  let showPopover = false;

  let modeOptions = [
    { label: 'Mode 1', value: 'm1' },
    { label: 'Mode 2', value: 'm2' },
    { label: 'Mode 3', value: 'm3' }
  ];
  let selectedMode = 'm1';

  let stepperSteps = [
    { label: 'Assess' },
    { label: 'Route' },
    { label: 'Evacuate' }
  ];

  let breadcrumbPaths = [
    { name: 'Dashboard', href: '#' },
    { name: 'Barpeta', href: '#' },
    { name: 'Sorbhog', href: '#' }
  ];

  let radioOptions = [
    { value: 'osm', label: 'OpenStreetMap' },
    { value: 'satellite', label: 'Satellite (Bhuvan)' }
  ];
  let selectedLayer = 'satellite';

  let dropdownItems = [
    { text: 'View Full Report', action: () => alert('Viewing report') },
    { text: 'Export as CSV', action: () => alert('Exporting...') }
  ];

  let tableHeaders = [
    { key: 'name', label: 'Habitation', sortable: true },
    { key: 'district', label: 'District', sortable: true },
    { key: 'score', label: 'Risk Score', sortable: true }
  ];
  
  let tableData = [
    { name: 'Nongpoh', district: 'Ri Bhoi', score: 9.1 },
    { name: 'Sorbhog', district: 'Barpeta', score: 8.7 },
    { name: 'Bahari', district: 'Barpeta', score: 6.2 }
  ];
</script>

<main class="style-guide">
  <h1>Component Library</h1>
  <p>Testing the official government theme.</p>

  <section>
    <h2>Core Primitives</h2>
    <div style="max-width: 450px;">
      <Tabs tabs={['Overview', 'Risk Zones', 'Safe Sites']} />
      
      <div style="margin-top: 2rem;">
        <Switch label="Satellite Map View" checked={true} />
        <Checkbox label="Highlight High Risk Areas" checked={true} />
        <Checkbox label="Show Evacuation Routes" checked={false} />
      </div>

      <div style="margin-top: 2rem; display: flex; align-items: center; gap: 1rem;">
        <span>Loading Data: </span>
        <Spinner />
      </div>

      <Divider />

      <div style="display: flex; gap: 1rem; align-items: center;">
        <Button on:click={() => showModal = true}>Open Modal</Button>
        <Button on:click={() => showNotification = true}>Test Notification</Button>
        <Tooltip text="Sync with main database">
          <Button>Hover Me</Button>
        </Tooltip>
      </div>
    </div>
  </section>

  <Notification bind:visible={showNotification} message="Evacuation protocol initiated successfully." type="success" />

  <Modal bind:isOpen={showModal} title="Confirm Evacuation">
    <p>Are you sure you want to assign 3,200 people to Relocation Site A?</p>
    <div style="display: flex; gap: 1rem; margin-top: 2rem;">
      <Button on:click={() => showModal = false}>Confirm Action</Button>
      <Button on:click={() => showModal = false}>Cancel</Button>
    </div>
  </Modal>

  <section>
    <h2>Form Controls</h2>
    <div style="max-width: 450px;">
      <TextInput label="Resource Name" placeholder="e.g. Relief Tents" />
      
      <Select 
        label="Urgency Level" 
        options={[
          { value: 'low', label: 'Low - Monitor' },
          { value: 'medium', label: 'Medium - Prepare' },
          { value: 'high', label: 'High - Immediate Action' }
        ]} 
      />
      
      <Slider label="Evacuation Progress (%)" value={45} />
      
      <ProgressBar label="Site Capacity" value={82} />
    </div>
  </section>

  <section>
    <h2>Data Display</h2>
    <div style="max-width: 450px;">
      <Card title="Habitation: Sorbhog">
        <Badge tier="immediate" isBand slot="band">Immediate Risk</Badge>
        <p>District: Barpeta</p>
        <p>Population: 3,200</p>
        
        <Divider />
        
        <Accordion title="View Logistics Data">
          <p>Tents: 150/400 needed</p>
          <p>Food Rations: 3 Days remaining</p>
          <p>Medical Kits: Adequate</p>
        </Accordion>
      </Card>
    </div>
  </section>

  <Notification bind:visible={showNotification} message="Evacuation protocol initiated successfully." type="success" />

  <Modal bind:isOpen={showModal} title="Confirm Evacuation">
    <p>Are you sure you want to assign 3,200 people to Relocation Site A?</p>
    <div style="display: flex; gap: 1rem; margin-top: 2rem;">
      <Button on:click={() => showModal = false}>Confirm Action</Button>
      <Button on:click={() => showModal = false}>Cancel</Button>
    </div>
  </Modal>

  <section>
    <h2>Form Controls</h2>
    <div style="max-width: 450px;">
      <TextInput label="Resource Name" placeholder="e.g. Relief Tents" />
      
      <Select 
        label="Urgency Level" 
        options={[
          { value: 'low', label: 'Low - Monitor' },
          { value: 'medium', label: 'Medium - Prepare' },
          { value: 'high', label: 'High - Immediate Action' }
        ]} 
      />
      
      <Slider label="Evacuation Progress (%)" value={45} />
    </div>
  </section>

  <section>
    <h2>Buttons</h2>
    <div class="row">
      <Button>Primary Action</Button>
      <Button disabled>Disabled Action</Button>
    </div>
    <div class="row" style="margin-top: 1rem; max-width: 300px;">
      <Button fullWidth>Full Width Button</Button>
    </div>
  </section>

  <section>
    <h2>Risk Tiers (Badges)</h2>
    <div class="row">
      <Badge tier="immediate">Immediate</Badge>
      <Badge tier="moderate">Moderate</Badge>
      <Badge tier="safe">Safe / Inactive</Badge>
      <Badge tier="habitation">Habitation</Badge>
    </div>
  </section>

  <section>
    <h2>Data Cards</h2>
    <div style="max-width: 450px;">
      <Card title="Habitation: Sorbhog">
        <Badge tier="immediate" isBand slot="band">Immediate Risk</Badge>
        <p>District: Barpeta</p>
        <p>Population: 3,200</p>
      </Card>
      
      <Card title="Relocation Site">
        <p>Capacity: 1,200</p>
        <p>Suitability: 7.8/10</p>
        <div style="margin-top: 1rem;">
          <Button fullWidth>Assign Population</Button>
        </div>
      </Card>
    </div>
  </section>
  <section>
    <h2>New Data Primitives</h2>
    <div style="max-width: 700px; display: flex; flex-direction: column; gap: 2rem;">
      <div>
        <h4>Data Table</h4>
        <DataTable headers={tableHeaders} data={tableData} sortable={true} />
      </div>

      <div style="display: flex; gap: 3rem;">
        <div>
          <h4>Metrics</h4>
          <Metric label="Total Population" value="24,500" />
        </div>
        <div>
          <h4>Tags</h4>
          <div style="display: flex; gap: 0.5rem; margin-top: 0.5rem;">
            <Tag>Confirmed</Tag>
            <Tag text="Unconfirmed" />
          </div>
        </div>
      </div>

      <div>
        <h4>Callout (Action Center)</h4>
        <Callout type="info">
          Sorbhog is flagged Immediate priority, hazard score 8.7/10.
          Nearest viable relocation site: Kalgachia Uplands (capacity 850, 2.4km away).
        </Callout>
        <Callout type="alert">
          Warning: Cloudburst threshold exceeded. Triggering Mode 2 rerouting.
        </Callout>
      </div>
      
      <div>
        <h4>Skeleton Loading</h4>
        <div style="display: flex; flex-direction: column; gap: 0.5rem;">
          <Skeleton height="1.5rem" width="60%" />
          <Skeleton height="1rem" width="100%" />
          <Skeleton height="1rem" width="80%" />
        </div>
      </div>
    </div>
  </section>

  <section>
    <h2>New Form & Overlay Primitives</h2>
    <div style="max-width: 450px; display: flex; flex-direction: column; gap: 1.5rem;">
      <SearchInput label="Search Locations" />
      <Textarea label="Report Description" placeholder="Describe the road blockage..." />
      <FileInput label="Attach Evidence (Photo)" />
      
      <div>
        <Button on:click={() => showDrawer = true}>Open Data Provenance Drawer</Button>
      </div>
    </div>
  </section>

  <Drawer bind:isOpen={showDrawer} title="Data Provenance" position="right">
    <p><strong>Elevation Source:</strong> CartoDEM via Bhoonidhi</p>
    <p><strong>Landslide Inventory:</strong> GSI NLSM</p>
    <p><strong>Flood Hazard:</strong> Bhuvan Web Services</p>
    <div style="margin-top: 2rem;">
      <Button fullWidth on:click={() => showDrawer = false}>Close Settings</Button>
    </div>
  </Drawer>

  <section>
    <h2>Final Edge-Case Primitives</h2>
    <div style="max-width: 600px; display: flex; flex-direction: column; gap: 2rem;">
      <div>
        <h4>Breadcrumbs</h4>
        <Breadcrumbs paths={breadcrumbPaths} />
      </div>

      <div>
        <h4>Map Legend</h4>
        <MapLegend />
      </div>

      <div>
        <h4>Radio Group (Layer Selection)</h4>
        <RadioGroup options={radioOptions} bind:selected={selectedLayer} />
      </div>

      <div>
        <h4>Pagination</h4>
        <Pagination currentPage={3} totalPages={12} />
      </div>

      <div>
        <h4>Dropdown Menu</h4>
        <DropdownMenu label="Row Actions" items={dropdownItems} />
      </div>

      <div>
        <h4>Icons</h4>
        <div style="display: flex; gap: 1rem; color: var(--text);">
          <Icon name="search" size="24" />
          <Icon name="check" size="24" />
          <Icon name="x" size="24" />
          <Icon name="chevronDown" size="24" />
        </div>
      </div>
    </div>
  </section>

  <section>
    <h2>Dashboard View Primitives</h2>
    <div style="max-width: 600px; display: flex; flex-direction: column; gap: 2rem;">
      
      <div>
        <h4>Empty State</h4>
        <EmptyState 
          title="No Active Events" 
          description="Monitoring baseline hazard data. Waiting for threshold cross." 
          icon="check"
        />
      </div>

      <div>
        <h4>Button Group (Segmented Control)</h4>
        <ButtonGroup options={modeOptions} bind:selected={selectedMode} />
        <p style="margin-top: 0.5rem; font-size: 0.85rem; color: color-mix(in srgb, var(--text) 60%, transparent);">Selected: {selectedMode}</p>
      </div>

      <div>
        <h4>Popover</h4>
        <div style="position: relative;">
          <Button on:click={() => showPopover = !showPopover}>Toggle Popover</Button>
          <Popover bind:isOpen={showPopover}>
            <p><strong>Nongpoh Marker</strong></p>
            <p style="font-size: 0.85rem; margin-top: 0.5rem;">Population: 4,200</p>
          </Popover>
        </div>
      </div>

      <div>
        <h4>Image Block (Citizen Report)</h4>
        <ImageBlock 
          src="https://via.placeholder.com/600x300.png?text=Bridge+Collapse+Evidence"
          caption="Reported at 14:32 IST"
        />
      </div>

      <div>
        <h4>Stepper (Evacuation Flow)</h4>
        <Stepper steps={stepperSteps} currentStep={1} />
      </div>
      
    </div>
  </section>

  <section>
    <h2>Domain-Specific Primitives (Suraksha Setu)</h2>
    <div style="max-width: 700px; display: flex; flex-direction: column; gap: 2rem;">
      
      <div>
        <h4>Location Picker</h4>
        <LocationPicker label="Incident Coordinates" />
      </div>

      <div>
        <h4>Action Center Card</h4>
        <ActionCenterCard 
          habitationName="Sorbhog" 
          tier="immediate" 
          score="0.87" 
          siteName="Kalgachia Uplands" 
          capacity={850} 
          distance={2.4} 
        />
      </div>

      <div>
        <h4>Map Popup Preview</h4>
        <div style="border: 1px solid var(--secondary); padding: 1rem; border-radius: 4px; display: inline-block; background: var(--background);">
          <MapPopup 
            title="Nongpoh" 
            tier="immediate" 
            population={4200} 
            hazardScore={0.91} 
            onAction={() => alert('View Details clicked')} 
          />
        </div>
      </div>

      <div>
        <h4>Leaflet Map Primitive</h4>
        <Map />
      </div>

    </div>
  </section>

</main>

<style>
  .style-guide {
    padding: 3rem;
    overflow-y: auto;
    height: 100vh;
  }

  section {
    margin-top: 2rem;
    margin-bottom: 3rem;
    padding-bottom: 2rem;
    border-bottom: 1px solid var(--secondary);
  }

  p {
    margin-bottom: 0.5rem;
  }

  .row {
    display: flex;
    gap: 1rem;
    align-items: center;
    flex-wrap: wrap;
  }
</style>
