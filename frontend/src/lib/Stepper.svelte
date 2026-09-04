<script>
  export let steps = []; // Array of { label: string }
  export let currentStep = 0; // 0-indexed
</script>

<div class="gov-stepper">
  {#each steps as step, i}
    <div class="step" class:completed={i < currentStep} class:active={i === currentStep}>
      <div class="step-circle">{i + 1}</div>
      <div class="step-label">{step.label}</div>
    </div>
    {#if i < steps.length - 1}
      <div class="step-line" class:completed={i < currentStep}></div>
    {/if}
  {/each}
</div>

<style>
  .gov-stepper {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
  }

  .step {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    position: relative;
    z-index: 2;
  }

  .step-circle {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background-color: var(--background);
    border: 2px solid var(--secondary);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.85rem;
    font-weight: 400;
    color: color-mix(in srgb, var(--text) 60%, transparent);
    transition: all 0.2s;
  }

  .step-label {
    font-size: 0.75rem;
    color: color-mix(in srgb, var(--text) 60%, transparent);
    text-align: center;
  }

  .step.active .step-circle {
    border-color: var(--primary-dark);
    color: var(--text);
  }

  .step.active .step-label {
    color: var(--text);
  }

  .step.completed .step-circle {
    background-color: var(--primary-dark);
    border-color: var(--primary-dark);
    color: var(--background);
  }

  .step-line {
    flex: 1;
    height: 2px;
    background-color: var(--secondary);
    margin: 0 0.5rem;
    transform: translateY(-10px); /* aligns line visually with circles */
    transition: background-color 0.2s;
  }

  .step-line.completed {
    background-color: var(--primary-dark);
  }
</style>
