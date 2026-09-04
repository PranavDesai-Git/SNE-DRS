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

  let showModal = false;
  let showNotification = false;
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
