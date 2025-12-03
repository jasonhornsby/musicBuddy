<script lang="ts">
	import { getAudioContext } from "$lib/context/audio.svelte";
	import { Spinner } from "$lib/components/ui/spinner";
	import { onMount } from "svelte";


    const audioContext = getAudioContext();

    onMount(() => {
        if (audioContext.metadata) return;
        audioContext.getAudioMetadata();
    })
</script>

<header class="p-4">
	<h2 class="text-lg font-semibold">General Information</h2>
</header>

<div class="flex flex-col gap-2 p-4">
	{#if audioContext.metadata}
		<div class="flex flex-col gap-2">
			<p class="text-sm text-muted-foreground">
				Sample Rate: {audioContext.metadata.sampleRate}
			</p>
			<p class="text-sm text-muted-foreground">Channels: {audioContext.metadata.channels}</p>
			<p class="text-sm text-muted-foreground">
				Duration: {audioContext.metadata.durationMs}ms
			</p>
		</div>
	{:else}
		<Spinner class="size-10 mx-auto" />
	{/if}
</div>
