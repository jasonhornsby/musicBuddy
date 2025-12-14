<script lang="ts">
	import { getAudioContext } from "$lib/context/audio.svelte";
	import { WaveformVisualisation } from "$lib/utils/visualiser";

	const audioContext = getAudioContext();
	const workerManager = audioContext.getAudioWorkerManager();

	const width = 2000;
	const height = $state(200);
	const visualisation = new WaveformVisualisation(width);
	workerManager.createVisualizer(visualisation);

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

		visualisation.registerContext(ctx);

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

		visualisation.draw();
	}

	$effect(() => {
		if (canvas) {
			setupCanvas();
			draw();
		}
	});

	$effect(() => {
		if (!canvas) return;

		const resizeObserver = new ResizeObserver((entries) => {
			for (const entry of entries) {
				if (entry.target === canvas) {
					draw();
				}
			}
		});

		resizeObserver.observe(canvas);

		return () => {
			resizeObserver.disconnect();
		};
	});
</script>

<div class="flex flex-col gap-4">
	<canvas
		bind:this={canvas}
		style="height: {height}px;"
		class="rounded border border-neutral-700 w-full"
	></canvas>
</div>
