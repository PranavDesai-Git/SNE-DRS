<script>
  // Simple icon primitive wrapper
  export let name = "";
  export let size = "16";
  
  // Minimal internal dictionary for raw SVGs so we don't need external libraries
  // Can be expanded as needed.
  const icons = {
    search: '<circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line>',
    check: '<polyline points="20 6 9 17 4 12"></polyline>',
    x: '<line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line>',
    chevronDown: '<polyline points="6 9 12 15 18 9"></polyline>'
  };
</script>

<svg 
  class="gov-icon" 
  xmlns="http://www.w3.org/2000/svg" 
  width={size} 
  height={size} 
  viewBox="0 0 24 24" 
  fill="none" 
  stroke="currentColor" 
  stroke-width="2" 
  stroke-linecap="round" 
  stroke-linejoin="round"
>
  {@html icons[name] || ''}
</svg>

<style>
  .gov-icon {
    display: inline-block;
    vertical-align: middle;
  }
</style>
