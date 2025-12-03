<script lang="ts">
	import { getAudioContext } from "$lib/context/audio.svelte";
	import { Spinner } from "$lib/components/ui/spinner";
	import { Separator } from "$lib/components/ui/separator";
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
</script>

<div class="flex flex-col h-full">
	<header class="px-5 py-4 flex items-center justify-between">
		<h2 class="text-sm font-medium tracking-wide uppercase text-muted-foreground">Overview</h2>
	</header>

	<Separator />

	<div class="flex-1 p-5">
		{#if audioContext.metadata}
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-1.5 p-4 rounded-lg bg-muted/50">
					<p class="text-xs font-medium uppercase tracking-wider text-muted-foreground">
						Sample Rate
					</p>
					<p class="text-xl font-semibold tabular-nums">
						{formatSampleRate(audioContext.metadata.sampleRate)}
					</p>
				</div>

				<div class="space-y-1.5 p-4 rounded-lg bg-muted/50">
					<p class="text-xs font-medium uppercase tracking-wider text-muted-foreground">Channels</p>
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
					<p class="text-xs font-medium uppercase tracking-wider text-muted-foreground">Duration</p>
					<p class="text-xl font-semibold tabular-nums">
						{formatDuration(audioContext.metadata.durationMs)}
						<span class="text-sm font-normal text-muted-foreground ml-1">
							({audioContext.metadata.durationMs.toLocaleString()}ms)
						</span>
					</p>
				</div>
			</div>
		{:else}
			<div class="flex items-center justify-center h-32">
				<Spinner class="size-8" />
			</div>
		{/if}
	</div>
</div>
