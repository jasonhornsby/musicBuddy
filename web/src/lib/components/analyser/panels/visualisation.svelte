<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import { getAudioContext } from "$lib/context/audio.svelte";

	const audioContext = getAudioContext();
	const workerManager = audioContext.getAudioWorkerManager();

	const width = 800;
	const height = $state(300);
	const vizData = workerManager.createVisualizer('test-viz', 'waveform', width);

	let canvas = $state<HTMLCanvasElement | null>(null);
	let dpr = $state(typeof window !== 'undefined' ? window.devicePixelRatio : 1);

	function setupCanvas() {
		if (!canvas) return null;

		const ctx = canvas.getContext('2d');
		if (!ctx) return null;

		// Set the canvas buffer size scaled by DPI
		canvas.width = width * dpr;
		canvas.height = height * dpr;

		// Scale the context to account for DPI
		ctx.scale(dpr, dpr);

		return ctx;
	}

	function draw() {
		if (!canvas) return;

		const ctx = canvas.getContext('2d');
		if (!ctx) return;

		// Reset transform and reapply DPI scaling
		ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

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

	function updateViz() {
		workerManager.updateVisualizer('test-viz');
		draw();
	}

	$effect(() => {
		if (canvas) {
			setupCanvas();
			draw();
		}
	});
</script>

<div class="flex flex-col gap-4">
	<div class="flex gap-2">
		<Button onclick={updateViz}>UpdateViz</Button>
	</div>
	<canvas
		bind:this={canvas}
		style="width: {width}px; height: {height}px;"
		class="rounded border border-neutral-700"
	></canvas>
</div>
