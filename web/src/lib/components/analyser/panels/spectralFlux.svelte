<script lang="ts">
	import { getAudioContext } from "$lib/context/audio.svelte";
	import * as Card from "$lib/components/ui/card";

	const audioContext = getAudioContext();
</script>

<div class="flex flex-col h-full">
	<div class="flex-1 p-5 space-y-5 overflow-y-auto">
		<Card.Root class="border-dashed border-2 border-muted-foreground/30 bg-muted/20">
			<Card.Header class="text-center py-12">
				<div class="mx-auto mb-4 size-16 rounded-full bg-muted/50 flex items-center justify-center">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="size-8 text-muted-foreground/60"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.5"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d="M2 12h2" />
						<path d="M6 8v8" />
						<path d="M10 4v16" />
						<path d="M14 6v12" />
						<path d="M18 9v6" />
						<path d="M22 12h-2" />
					</svg>
				</div>
				<Card.Title class="text-lg text-muted-foreground">Spectral Flux Analysis</Card.Title>
				<Card.Description class="text-sm text-muted-foreground/70 max-w-xs mx-auto mt-2">
					Spectral flux visualization coming soon. This will show the rate of change in the power
					spectrum over time.
				</Card.Description>
			</Card.Header>
		</Card.Root>

		<div class="space-y-3">
			<h3 class="text-xs font-medium uppercase tracking-wider text-muted-foreground px-1">
				What is Spectral Flux?
			</h3>
			<p class="text-sm text-muted-foreground/80 leading-relaxed px-1">
				Spectral flux measures how quickly the power spectrum of a signal changes over time. It's
				commonly used for onset detection, beat tracking, and identifying sudden changes in the
				audio signal's frequency content.
			</p>
		</div>
	</div>
</div>
