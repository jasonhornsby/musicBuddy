<script lang="ts">
    import { onMount } from 'svelte';
    import {Button} from '$lib/components/ui/button';
	import { Spinner } from '$lib/components/ui/spinner';
	import type { DemoFile } from '../remote/demo-files.remote';
	import { Separator } from '$lib/components/ui/separator';
	import { Analyser } from '$lib/components/analyser';

    const { data } = $props();

    let isWasmLoaded = $state(false);
    let loadingAudio = $state(false);
    let audioLoaded = $state(false);
    let uploadFileInput = $state<HTMLInputElement | null>(null);
  
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
      
      isWasmLoaded = true;
      console.log("Go Wasm loaded via SvelteKit!");
    });

    async function handleFileChange(event: Event) {
        loadingAudio = true
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
        loadingAudio = false;   
        audioLoaded = true;
    }

    async function loadDemoFile(demoFile: DemoFile) {
        loadingAudio = true;
        const response = await fetch(demoFile.src);

        const arrayBuffer = await response.arrayBuffer();
        const uint8Array = new Uint8Array(arrayBuffer);
        const result = await window.loadAudio(uint8Array);
        console.log("result: ", result);
        loadingAudio = false;
        audioLoaded = true;
    }
</script>

{#if isWasmLoaded}
	{#if !audioLoaded}
		{#if loadingAudio}
			<div class="h-full w-full flex items-center justify-center">
				<div class="flex flex-col items-center justify-center">
					<Spinner class="size-10 text-muted-foreground" />
					<span class="text-sm text-muted-foreground">Parsing audio</span>
				</div>
			</div>
		{:else}
			<div class="p-4 flex flex-col gap-4">
				<h1 class="text-2xl font-bold">Analyse your music</h1>
				<h2>Demo files</h2>
				<div class="flex flex-row items-center gap-4">
					{#each data.demoFiles as demoFile}
						<Button variant="default" onclick={() => loadDemoFile(demoFile)}>{demoFile.name}</Button
						>
					{/each}
				</div>
				<Separator />
				<h2>Upload your own file</h2>

				<div class="self-start">
					<Button onclick={() => uploadFileInput?.click()}>Upload file</Button>
					<input
						bind:this={uploadFileInput}
						type="file"
						class="hidden"
						onchange={handleFileChange}
						accept="audio/mp3"
					/>
				</div>
			</div>
		{/if}
	{:else}
		<Analyser />
	{/if}
{:else}
	<div class="h-dvh w-dvw flex items-center justify-center">
		<div class="flex flex-col items-center justify-center">
			<Spinner class="size-10 text-muted-foreground" />
			<span class="text-sm text-muted-foreground">Loading</span>
		</div>
	</div>
{/if}
