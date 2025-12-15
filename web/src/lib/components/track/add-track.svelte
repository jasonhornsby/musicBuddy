<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button';
	import { Spinner } from '$lib/components/ui/spinner';
	import { Separator } from '$lib/components/ui/separator';
	import { CloudUpload, PlayCircle } from 'lucide-svelte';
	import { getAudioContext } from '$lib/context/audio.svelte';

	interface DemoFile {
		name: string;
		src: string;
	}

	const demoFiles: DemoFile[] = [
		{ name: 'Jazz', src: 'examples/jazz.mp3' },
		{ name: 'Hip Hop', src: 'examples/hiphop.mp3' },
		{ name: 'Error', src: 'examples/error.mp3' }
	];

	let { open = $bindable(false) } = $props();

	let uploadFileInput = $state<HTMLInputElement | null>(null);
	let isDragging = $state(false);

	const audioContext = getAudioContext();

	async function handleFileChange(event: Event) {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		await audioContext.loadAudioFile(file);
		open = false;
	}

	async function loadDemoFile(demoFile: DemoFile) {
		await audioContext.loadAudioFromSrc(demoFile.src);
		open = false;
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
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="sm:max-w-sm">
		<Dialog.Header>
			<Dialog.Title>Add Track</Dialog.Title>
			<Dialog.Description>Upload an audio file or try a demo track</Dialog.Description>
		</Dialog.Header>

		<div class="flex flex-col gap-4">
			{#if audioContext.isParsingAudio}
				<div class="flex flex-col items-center justify-center py-8 gap-2">
					<Spinner class="size-6 text-primary" />
					<span class="text-sm text-muted-foreground">Processing audio...</span>
				</div>
			{:else}
				<button
					type="button"
					class="group flex flex-col items-center justify-center gap-2 rounded-lg border-2 border-dashed border-border/60 bg-muted/30 p-5 transition-all hover:border-primary/50 hover:bg-muted/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring {isDragging
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
					<CloudUpload
						class="h-5 w-5 text-muted-foreground group-hover:text-primary transition-colors"
					/>
					<div class="text-center">
						<p class="text-sm font-medium text-foreground">Drop your audio file here</p>
						<p class="text-xs text-muted-foreground">or click to browse • MP3 files supported</p>
					</div>
					<input
						bind:this={uploadFileInput}
						type="file"
						class="hidden"
						onchange={handleFileChange}
						accept="audio/mp3"
					/>
				</button>

				<div class="flex items-center gap-3">
					<Separator class="flex-1" />
					<span class="text-xs text-muted-foreground">or try a demo</span>
					<Separator class="flex-1" />
				</div>

				<div class="grid grid-cols-3 gap-2">
					{#each demoFiles as demoFile}
						<Button
							variant="outline"
							size="sm"
							class="h-auto py-2"
							onclick={() => loadDemoFile(demoFile)}
						>
							<PlayCircle class="h-4 w-4 mr-1" strokeWidth={1.5} />
							{demoFile.name}
						</Button>
					{/each}
				</div>
			{/if}
		</div>
	</Dialog.Content>
</Dialog.Root>
