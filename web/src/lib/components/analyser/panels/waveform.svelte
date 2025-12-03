<script lang="ts">
	import { getAudioContext } from "$lib/context/audio.svelte";
	import WaveformPlotter from "./waveformPlotter.svelte";

    const audioContext = getAudioContext();
</script>

<div class="h-full w-full flex-1">
	{#if audioContext.decoded}
		<WaveformPlotter />
	{/if}
</div>
