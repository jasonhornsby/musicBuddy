<script lang="ts">
	import { getAudioContext } from "$lib/context/audio.svelte";
	import { WaveformVisualisation } from "$lib/utils/visualiser";

	const audioContext = getAudioContext();
	const workerManager = audioContext.getAudioWorkerManager();

	const visualisation = new WaveformVisualisation(300);
	
	workerManager.createVisualizer(visualisation);

	let canvas = $state<HTMLCanvasElement | null>(null);
	let dpr = $state(typeof window !== 'undefined' ? window.devicePixelRatio : 1);
	let canvasWidth = $state(0);
	let canvasHeight = $state(0);

	function setupCanvas() {
		if (!canvas) return null;

		const ctx = canvas.getContext('2d');
		if (!ctx) return null;

	const rect = canvas.getBoundingClientRect();
	canvasWidth = Math.max(1, Math.floor(rect.width));
	canvasHeight = Math.max(1, Math.floor(rect.height));

	// Match the backing buffer to the on-screen size to prevent stretching.
	const scaledWidth = canvasWidth * dpr;
	const scaledHeight = canvasHeight * dpr;
	if (canvas.width !== scaledWidth || canvas.height !== scaledHeight) {
		canvas.width = scaledWidth;
		canvas.height = scaledHeight;
	}

	// Ensure the context accounts for device pixel ratio.
	ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

		visualisation.registerContext(ctx);

		return ctx;
	}

	function draw() {
		if (!canvas) return;

		const ctx = canvas.getContext('2d');
		if (!ctx) return;

	// Ensure the buffer matches the current element size (handles resizes).
	setupCanvas();

		// Clear canvas (matching track.svelte white background)
		ctx.fillStyle = '#ffffff';
		ctx.fillRect(0, 0, canvasWidth, canvasHeight);

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
					setupCanvas();
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

<canvas bind:this={canvas} class="w-full h-full"></canvas>
