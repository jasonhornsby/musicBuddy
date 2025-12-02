<script lang="ts">
    import { onMount } from 'svelte';
    import {Button} from '$lib/components/ui/button';
	import { Spinner } from '$lib/components/ui/spinner';
  
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

    async function handleFileChange(event: Event) {
        const input = event.target as HTMLInputElement;
        const file = input.files?.[0];
        if (!file) return;

        try {
            const arrayBuffer = await file.arrayBuffer();
            const uint8Array = new Uint8Array(arrayBuffer);
            const result = await window.loadAudio(uint8Array);
            console.log("result: ", result);
        } catch (e) {
            console.error("Error loading audio: ", e);
        }
    }
</script>

{#if isLoaded}
	<h1>SvelteKit + Go Wasm</h1>
	<input type="file" on:change={handleFileChange} accept="audio/mp3" />
{:else}
	<div class="h-dvh w-dvw flex items-center justify-center">
		<div class="flex flex-col items-center justify-center">
			<Spinner class="size-10 text-muted-foreground" />
			<span class="text-sm text-muted-foreground">Loading</span>
		</div>
	</div>
{/if}
