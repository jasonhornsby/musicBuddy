<script lang="ts">
	import { getAudioContext } from "$lib/context/audio.svelte";
	import Plotter from "./plotter.svelte";

    const audioContext = getAudioContext();
</script>

<div class="h-full w-full bg-red-50 flex-1">
	{#if audioContext.decoded}
		<Plotter />
	{/if}
</div>
