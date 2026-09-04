<script>
  export let currentPage = 1;
  export let totalPages = 1;

  function prev() {
    if (currentPage > 1) currentPage -= 1;
  }

  function next() {
    if (currentPage < totalPages) currentPage += 1;
  }
</script>

<div class="gov-pagination">
  <button class="page-btn" on:click={prev} disabled={currentPage === 1}>
    &lt;
  </button>
  <span class="page-info">Page {currentPage} of {totalPages}</span>
  <button class="page-btn" on:click={next} disabled={currentPage === totalPages}>
    &gt;
  </button>
</div>

<style>
  .gov-pagination {
    display: inline-flex;
    align-items: center;
    gap: 1rem;
    font-size: 0.9rem;
  }

  .page-info {
    color: var(--text);
  }

  .page-btn {
    background-color: var(--background);
    border: 1px solid var(--secondary);
    color: var(--text);
    padding: 0.25rem 0.75rem;
    border-radius: 4px;
    cursor: pointer;
    font-family: 'Poppins', sans-serif;
    transition: background-color 0.2s, border-color 0.2s;
  }

  .page-btn:hover:not(:disabled) {
    border-color: var(--primary-dark);
    background-color: color-mix(in srgb, var(--background), black 2%);
  }

  .page-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
