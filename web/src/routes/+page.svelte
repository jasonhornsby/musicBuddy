<script>
    import { onMount } from 'svelte';
    import {Button} from '$lib/components/ui/button';
  
    let result = "...";
    let isLoaded = false;
  
    onMount(async () => {
      // 1. Initialize the Go wrapper
      const go = new window.Go();
      
      // 2. Fetch the symlinked wasm file from /static
      // Note: fetch("/main.wasm") works because 'static' is the root 
      const response = await fetch("/main.wasm");
      const buffer = await response.arrayBuffer();
      
      // 3. Instantiate
      const { instance } = await WebAssembly.instantiate(buffer, go.importObject);
      
      // 4. Run the Go program (this registers window.add)
      go.run(instance);
      
      isLoaded = true;
      console.log("Go Wasm loaded via SvelteKit!");
    });
  
    function handleCalculate() {
      if (!isLoaded) return;
      // Call the function exposed by Go
      result = window.add(10, 20); 
      console.log("result: ", result);
    }
</script>

<h1>SvelteKit + Go Wasm</h1>

<Button onclick={handleCalculate} disabled={!isLoaded}>Calculate 10 + 20</Button>

<p>Result: {result}</p>
