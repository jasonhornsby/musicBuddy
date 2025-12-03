<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Spinner } from '$lib/components/ui/spinner';
	import type { DemoFile } from '../remote/demo-files.remote';
	import { Separator } from '$lib/components/ui/separator';
	import { Analyser } from '$lib/components/analyser';
	import * as Card from '$lib/components/ui/card';
	import { Music, CloudUpload, PlayCircle } from 'lucide-svelte';
	import { getAudioContext, setAudioContext } from '$lib/context/audio.svelte.js';

	const { data } = $props();

	let isWasmLoaded = $state(false);
	let uploadFileInput = $state<HTMLInputElement | null>(null);
	let isDragging = $state(false);

    const audioContext = setAudioContext();

	onMount(async () => {
		const go = new window.Go();
		const response = await fetch('/main.wasm');
		const buffer = await response.arrayBuffer();
		const { instance } = await WebAssembly.instantiate(buffer, go.importObject);
		go.run(instance);
		isWasmLoaded = true;
	});

	async function handleFileChange(event: Event) {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;

		audioContext.loadAudio(file);
	}

	async function loadDemoFile(demoFile: DemoFile) {
		audioContext.loadAudioFromSrc(demoFile.src);
	}

	function handleDrop(event: DragEvent) {
		event.preventDefault();
		isDragging = false;
		const file = event.dataTransfer?.files[0];
		if (file && file.type === 'audio/mpeg') {
			const fakeEvent = { target: { files: [file] } } as unknown as Event;
			handleFileChange(fakeEvent);
		}
	}

    $inspect(audioContext.audioLoaded);
    $inspect(audioContext.parsingAudio);
    $inspect(isWasmLoaded);
</script>

{#if !audioContext.audioLoaded}
	<div
		class="min-h-screen w-full flex items-center justify-center p-6 bg-linear-to-br from-background via-background to-muted/30"
	>
		<Card.Root class="w-full max-w-md shadow-xl border-border/50">
			<Card.Header class="text-center pb-2">
				<div
					class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-primary/10"
				>
					<Music class="h-7 w-7 text-primary" strokeWidth={1.5} />
				</div>
				<Card.Title class="text-2xl font-semibold tracking-tight">Audio Analyser</Card.Title>
				<Card.Description class="text-muted-foreground">
					Discover the patterns and structure in your music
				</Card.Description>
			</Card.Header>

			<Card.Content class="flex flex-col gap-6">
				{#if audioContext.parsingAudio || !isWasmLoaded}
					<div class="flex flex-col items-center justify-center py-12 gap-3">
						<Spinner class="size-8 text-primary" />
						<span class="text-sm text-muted-foreground font-medium">
							{!isWasmLoaded ? 'Initializing...' : 'Processing audio...'}
						</span>
					</div>
				{:else}
					<!-- Upload Zone -->
					<button
						type="button"
						class="group relative flex flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed border-border/60 bg-muted/30 p-8 transition-all hover:border-primary/50 hover:bg-muted/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring {isDragging
							? 'border-primary bg-primary/5'
							: ''}"
						onclick={() => uploadFileInput?.click()}
						ondragover={(e) => {
							e.preventDefault();
							isDragging = true;
						}}
						ondragleave={() => (isDragging = false)}
						ondrop={handleDrop}
					>
						<div
							class="flex h-12 w-12 items-center justify-center rounded-full bg-background shadow-sm border border-border/50 group-hover:border-primary/30 transition-colors"
						>
							<CloudUpload
								class="h-5 w-5 text-muted-foreground group-hover:text-primary transition-colors"
							/>
						</div>
						<div class="text-center">
							<p class="text-sm font-medium text-foreground">Drop your audio file here</p>
							<p class="mt-1 text-xs text-muted-foreground">
								or click to browse • MP3 files supported
							</p>
						</div>
						<input
							bind:this={uploadFileInput}
							type="file"
							class="hidden"
							onchange={handleFileChange}
							accept="audio/mp3"
						/>
					</button>

					<div class="flex items-center gap-4">
						<Separator class="flex-1" />
						<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider"
							>or try a demo</span
						>
						<Separator class="flex-1" />
					</div>

					<!-- Demo Files -->
					<div class="grid grid-cols-2 gap-3">
						{#each data.demoFiles as demoFile}
							<Button
								variant="outline"
								class="h-auto py-4 flex flex-col gap-1 hover:bg-accent hover:border-primary/30"
								onclick={() => loadDemoFile(demoFile)}
							>
								<PlayCircle class="h-5 w-5 text-muted-foreground mb-1" strokeWidth={1.5} />
								<span class="font-medium">{demoFile.name}</span>
							</Button>
						{/each}
					</div>
				{/if}
			</Card.Content>
		</Card.Root>
	</div>
{:else}
	<Analyser />
{/if}
