<script>
  export let width = "100%";
  export let height = "1rem";
  export let borderRadius = "4px";
</script>

<div class="gov-skeleton" style="width: {width}; height: {height}; border-radius: {borderRadius};"></div>

<style>
  .gov-skeleton {
    background: linear-gradient(90deg, 
      var(--secondary) 25%, 
      color-mix(in srgb, var(--secondary), white 30%) 50%, 
      var(--secondary) 75%
    );
    background-size: 200% 100%;
    animation: loading 1.5s infinite;
    opacity: 0.3;
  }

  @keyframes loading {
    0% { background-position: 200% 0; }
    100% { background-position: -200% 0; }
  }
</style>
