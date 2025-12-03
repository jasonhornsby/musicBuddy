<script lang="ts">
	import { getAudioContext } from "$lib/context/audio.svelte";
	import { Spinner } from "$lib/components/ui/spinner";
	import { Separator } from "$lib/components/ui/separator";
	import * as Card from "$lib/components/ui/card";
	import { onMount } from "svelte";

	const audioContext = getAudioContext();

	onMount(() => {
		if (audioContext.metadata) return;
		audioContext.getAudioMetadata();
	});

	function formatDuration(ms: number): string {
		const seconds = Math.floor(ms / 1000);
		const minutes = Math.floor(seconds / 60);
		const remainingSeconds = seconds % 60;
		return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
	}

	function formatSampleRate(rate: number): string {
		return `${(rate / 1000).toFixed(1)} kHz`;
	}

	function formatBitrate(bitrate: number): string {
		return `${Math.round(bitrate / 1000)} kbps`;
	}

	function hasTrackInfo(metadata: typeof audioContext.metadata): boolean {
		if (!metadata?.metadata) return false;
		const { name, artist, album, year } = metadata.metadata;
		return !!(name || artist || album || year);
	}
</script>

<div class="flex flex-col h-full">
	<header class="px-5 py-4 flex items-center justify-between">
		<h2 class="text-sm font-medium tracking-wide uppercase text-muted-foreground">Overview</h2>
	</header>

	<Separator />

	<div class="flex-1 p-5 space-y-5 overflow-y-auto">
		{#if audioContext.metadata}
			<!-- Track Metadata Card -->
			{#if hasTrackInfo(audioContext.metadata)}
				{@const track = audioContext.metadata.metadata}
				<Card.Root class="border-primary/20 bg-primary/5 py-3 gap-3">
					<Card.Header class="px-4 pb-0">
						<div class="flex items-center gap-3">
							<div
								class="shrink-0 size-10 rounded-md bg-primary/15 flex items-center justify-center"
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="size-5 text-primary"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="1.5"
									stroke-linecap="round"
									stroke-linejoin="round"
								>
									<path d="M9 18V5l12-2v13" />
									<circle cx="6" cy="18" r="3" />
									<circle cx="18" cy="16" r="3" />
								</svg>
							</div>
							<div class="flex-1 min-w-0">
								{#if track.name}
									<Card.Title class="text-base font-semibold truncate">{track.name}</Card.Title>
								{:else}
									<Card.Title class="text-base text-muted-foreground italic"
										>Unknown Track</Card.Title
									>
								{/if}
								{#if track.artist}
									<p class="text-sm text-muted-foreground truncate">{track.artist}</p>
								{/if}
							</div>
						</div>
					</Card.Header>

					{#if track.album || track.year || track.format}
						<Card.Content class="px-4 pt-0">
							<div class="flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
								{#if track.album}
									<span class="flex items-center gap-1">
										<svg
											xmlns="http://www.w3.org/2000/svg"
											class="size-3"
											viewBox="0 0 24 24"
											fill="none"
											stroke="currentColor"
											stroke-width="2"
											stroke-linecap="round"
											stroke-linejoin="round"
										>
											<circle cx="12" cy="12" r="10" />
											<circle cx="12" cy="12" r="3" />
										</svg>
										{track.album}
									</span>
								{/if}
								{#if track.year}
									<span><span class="text-muted-foreground/60">Year:</span> {track.year}</span>
								{/if}
								{#if track.format}
									<span
										><span class="text-muted-foreground/60">Format:</span>
										<span class="uppercase">{track.format}</span></span
									>
								{/if}
							</div>
						</Card.Content>
					{/if}
				</Card.Root>
			{/if}

			<!-- Technical Details -->
			<div class="space-y-3">
				<h3 class="text-xs font-medium uppercase tracking-wider text-muted-foreground px-1">
					Technical Details
				</h3>
				<div class="grid grid-cols-2 gap-3">
					<div class="space-y-1.5 p-4 rounded-lg bg-muted/50">
						<p class="text-xs font-medium uppercase tracking-wider text-muted-foreground">
							Sample Rate
						</p>
						<p class="text-xl font-semibold tabular-nums">
							{formatSampleRate(audioContext.metadata.sampleRate)}
						</p>
					</div>

					<div class="space-y-1.5 p-4 rounded-lg bg-muted/50">
						<p class="text-xs font-medium uppercase tracking-wider text-muted-foreground">
							Channels
						</p>
						<p class="text-xl font-semibold tabular-nums">
							{audioContext.metadata.channels}
						</p>
					</div>

					<div class="space-y-1.5 p-4 rounded-lg bg-muted/50">
						<p class="text-xs font-medium uppercase tracking-wider text-muted-foreground">
							Decoded Bitrate
						</p>
						<p class="text-xl font-semibold tabular-nums">
							{formatBitrate(audioContext.metadata.decodedBitrate)}
						</p>
					</div>

					<div class="space-y-1.5 p-4 rounded-lg bg-muted/50">
						<p class="text-xs font-medium uppercase tracking-wider text-muted-foreground">
							Duration
						</p>
						<p class="text-xl font-semibold tabular-nums">
							{formatDuration(audioContext.metadata.durationMs)}
							<span class="text-sm font-normal text-muted-foreground ml-1">
								({audioContext.metadata.durationMs.toLocaleString()}ms)
							</span>
						</p>
					</div>
				</div>
			</div>
		{:else}
			<div class="flex items-center justify-center h-32">
				<Spinner class="size-8" />
			</div>
		{/if}
	</div>
</div>
