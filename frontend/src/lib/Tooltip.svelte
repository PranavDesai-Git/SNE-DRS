<script>
  export let text = "";
</script>

<div class="tooltip-wrapper">
  <slot />
  <div class="tooltip">{text}</div>
</div>

<style>
  .tooltip-wrapper {
    position: relative;
    display: inline-block;
  }

  .tooltip {
    visibility: hidden;
    background-color: var(--text);
    color: var(--background); /* Inverts colors for high contrast pop-out */
    text-align: center;
    border-radius: 4px;
    padding: 0.25rem 0.5rem;
    font-family: 'Poppins', sans-serif;
    font-size: 0.75rem;
    position: absolute;
    z-index: 1000;
    bottom: 125%; 
    left: 50%;
    transform: translateX(-50%);
    opacity: 0;
    transition: opacity 0.2s;
    white-space: nowrap;
    pointer-events: none;
    box-shadow: 0 4px 6px rgba(0,0,0,0.2);
  }

  .tooltip::after {
    content: "";
    position: absolute;
    top: 100%;
    left: 50%;
    margin-left: -5px;
    border-width: 5px;
    border-style: solid;
    border-color: var(--text) transparent transparent transparent;
  }

  .tooltip-wrapper:hover .tooltip {
    visibility: visible;
    opacity: 1;
  }
</style>
