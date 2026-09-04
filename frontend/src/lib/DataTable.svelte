<script>
  export let headers = [];
  export let data = [];
  export let sortable = false;
  export let onRowClick = null;
  export let selectedId = null;

  let sortColumn = null;
  let sortAscending = true;

  function handleSort(header) {
    if (!sortable || !header.sortable) return;
    
    if (sortColumn === header.key) {
      sortAscending = !sortAscending;
    } else {
      sortColumn = header.key;
      sortAscending = true;
    }
  }

  $: sortedData = sortColumn 
    ? [...data].sort((a, b) => {
        let valA = a[sortColumn];
        let valB = b[sortColumn];
        if (valA < valB) return sortAscending ? -1 : 1;
        if (valA > valB) return sortAscending ? 1 : -1;
        return 0;
      })
    : data;
</script>

<div class="table-container">
  <table class="gov-table">
    <thead>
      <tr>
        {#each headers as header}
          <th 
            on:click={() => handleSort(header)}
            class:sortable={sortable && header.sortable}
          >
            {header.label}
            {#if sortable && header.sortable && sortColumn === header.key}
              <span class="sort-icon">{sortAscending ? '▲' : '▼'}</span>
            {/if}
          </th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each sortedData as row}
        <tr 
          on:click={() => onRowClick && onRowClick(row)}
          class:selected={selectedId === row.id}
          style={onRowClick ? "cursor: pointer;" : ""}
        >
          {#each headers as header}
            <td>
              <slot name="cell" {row} {header}>
                {row[header.key]}
              </slot>
            </td>
          {/each}
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  .table-container {
    width: 100%;
    overflow-x: auto;
    overflow-y: auto;
    border: 1px solid var(--secondary);
    border-radius: 4px;
    background-color: var(--background);
    flex-shrink: 0;
    max-height: 400px;
  }

  .gov-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.9rem;
  }

  th, td {
    padding: 0.75rem 1rem;
    text-align: left;
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  }

  th {
    font-weight: 400; /* No bold text per design system */
    color: color-mix(in srgb, var(--text) 60%, transparent);
    background-color: color-mix(in srgb, var(--background), black 2%);
    position: relative;
  }

  .sortable {
    cursor: pointer;
  }

  .sortable:hover {
    background-color: color-mix(in srgb, var(--background), black 4%);
  }

  .sort-icon {
    font-size: 0.7rem;
    margin-left: 0.25rem;
  }

  tbody tr:hover {
    background-color: color-mix(in srgb, var(--background), var(--primary) 5%);
  }

  tbody tr:last-child td {
    border-bottom: none;
  }

  tr.selected {
    background-color: color-mix(in srgb, var(--background), var(--primary) 15%) !important;
  }
</style>
