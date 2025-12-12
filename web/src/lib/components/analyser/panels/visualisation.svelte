<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import { getAudioContext } from "$lib/context/audio.svelte";

	const audioContext = getAudioContext();
	const workerManager = audioContext.getAudioWorkerManager();

	const width = 800;
	const height = $state(300);
	const vizData = workerManager.createVisualizer('test-viz', 'waveform', width);

	let canvas = $state<HTMLCanvasElement | null>(null);
	let animationId = $state<number | null>(null);

	function draw() {
		if (!canvas) return;

		const ctx = canvas.getContext('2d');
		if (!ctx) return;

		// Clear canvas
		ctx.fillStyle = '#1a1a1a';
		ctx.fillRect(0, 0, width, height);

		// Draw min/max bars (data is interleaved: [min, max, min, max, ...])
		ctx.strokeStyle = '#00ff88';
		ctx.lineWidth = 1;

		const numBars = vizData.length / 2;
		const barWidth = width / numBars;

		for (let i = 0; i < numBars; i++) {
			const minVal = vizData[i * 2];
			const maxVal = vizData[i * 2 + 1];

			// Normalize values (assuming -1 to 1 range) to canvas coordinates
			const yMin = ((1 - minVal) / 2) * height;
			const yMax = ((1 - maxVal) / 2) * height;

			const x = i * barWidth + barWidth / 2;

			ctx.beginPath();
			ctx.moveTo(x, yMin);
			ctx.lineTo(x, yMax);
			ctx.stroke();
		}
	}

	function loop() {
		workerManager.updateVisualizer('test-viz');
		draw();
	}

	function startLoop() {
		if (animationId === null) {
			loop();
		}
	}

	function stopLoop() {
		if (animationId !== null) {
			cancelAnimationFrame(animationId);
			animationId = null;
		}
	}

	function updateViz() {
		workerManager.updateVisualizer('test-viz');
		draw();
	}
</script>

<div class="flex flex-col gap-4">
	<div class="flex gap-2">
		<Button onclick={updateViz}>UpdateViz</Button>
		<Button onclick={startLoop}>Start</Button>
		<Button onclick={stopLoop}>Stop</Button>
	</div>
	<canvas bind:this={canvas} {width} {height} class="rounded border border-neutral-700"></canvas>
</div>
