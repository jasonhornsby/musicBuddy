<script>
	import { Button } from "$lib/components/ui/button";
	import { getAudioContext } from "$lib/context/audio.svelte";


	const audioContext = getAudioContext();

	const workerManager = audioContext.getAudioWorkerManager();

	const width = $state(800);
	const vizData = workerManager.createVisualizer('test-viz', 'waveform', width);

	console.log(vizData);

	function updateViz() {
		workerManager.updateVisualizer('test-viz');
		console.log(vizData);
	}


</script>

<Button onclick={updateViz}>UpdateViz</Button>

<!--
<MainPlotter />
-->
