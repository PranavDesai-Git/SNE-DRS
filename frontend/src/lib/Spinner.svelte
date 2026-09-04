<script>
  export let size = "24px";
  export let color = "var(--primary-dark)";
</script>

<div class="spinner" style="width: {size}; height: {size}; border-top-color: {color};"></div>

<style>
  .spinner {
    border: 3px solid var(--secondary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
</style>
