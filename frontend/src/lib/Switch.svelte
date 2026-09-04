<script>
  export let label = "";
  export let checked = false;
</script>

<label class="switch-container">
  <span class="label-text">{label}</span>
  <div class="switch" class:checked>
    <input type="checkbox" bind:checked />
    <span class="slider"></span>
  </div>
</label>

<style>
  .switch-container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    cursor: pointer;
    margin-bottom: 0.5rem;
  }
  
  .label-text {
    font-size: 0.85rem;
    font-family: 'Poppins', sans-serif;
    color: var(--text);
  }

  .switch {
    position: relative;
    width: 40px;
    height: 22px;
  }

  .switch input {
    opacity: 0;
    width: 0;
    height: 0;
  }

  .slider {
    position: absolute;
    cursor: pointer;
    top: 0; left: 0; right: 0; bottom: 0;
    background-color: var(--secondary);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border-radius: 22px;
  }

  .slider:before {
    position: absolute;
    content: "";
    height: 16px;
    width: 16px;
    left: 3px;
    bottom: 3px;
    background-color: var(--background);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    border-radius: 50%;
  }

  .switch.checked .slider {
    background-color: var(--primary-dark);
  }

  .switch.checked .slider:before {
    transform: translateX(18px);
  }

  .switch-container:hover .switch:not(.checked) .slider:before {
    transform: scale(1.15);
  }

  .switch-container:active .switch:not(.checked) .slider:before {
    transform: scale(0.9);
  }

  .switch-container:hover .switch.checked .slider:before {
    transform: translateX(18px) scale(1.15);
  }

  .switch-container:active .switch.checked .slider:before {
    transform: translateX(18px) scale(0.9);
  }
</style>
