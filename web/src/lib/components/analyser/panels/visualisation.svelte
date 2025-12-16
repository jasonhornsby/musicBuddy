<script lang="ts">
	import { getAudioContext } from "$lib/context/audio.svelte";
	import { Visualisation, WaveformVisualisation } from "$lib/utils/visualiser.svelte";
	import type { ChannelMode } from "$lib/types/nodeConfig";

	const {visualisation}: {visualisation: Visualisation } = $props()

	const audioContext = getAudioContext();
	const workerManager = audioContext.getAudioWorkerManager();

	// Width of the left-hand legend/scale area in the waveform visualisation
	// Keep this in sync with `scaleTotalWidth` in `WaveformVisualisation.draw`
	const LEGEND_WIDTH = 24;

	let canvas = $state<HTMLCanvasElement | null>(null);
	let dpr = $state(typeof window !== 'undefined' ? window.devicePixelRatio : 1);
	let canvasWidth = $state(0);
	let canvasHeight = $state(0);
	let hoverX = $state<number | null>(null);

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
		console.log('drawing');

		const ctx = canvas.getContext('2d');
		if (!ctx) return;

		// Ensure the buffer matches the current element size (handles resizes).
		setupCanvas();

		// Clear canvas (matching track.svelte white background)
		ctx.fillStyle = '#ffffff';
		ctx.fillRect(0, 0, canvasWidth, canvasHeight);

		visualisation.draw();

		// Draw subtle vertical hover bar from top to bottom
		if (hoverX !== null) {
			ctx.save();
			ctx.strokeStyle = '#9ca3af'; // neutral-400 for subtle line
			ctx.globalAlpha = 0.4;
			ctx.lineWidth = 1;
			ctx.setLineDash([3, 3]);
			ctx.beginPath();
			ctx.moveTo(hoverX + 0.5, 0);
			ctx.lineTo(hoverX + 0.5, canvasHeight);
			ctx.stroke();
			ctx.restore();
		}
	}

	function handlePointerMove(event: PointerEvent) {
		if (!canvas) return;

		const rect = canvas.getBoundingClientRect();
		const x = event.clientX - rect.left;

		// If hovering over the legend area, don't show the hover line
		if (x < LEGEND_WIDTH) {
			hoverX = null;
			draw();
			return;
		}

		// Clamp the hover line to only appear over the waveform area (not the legend)
		hoverX = Math.max(LEGEND_WIDTH, Math.min(canvasWidth, x));
		draw();
	}

	function handlePointerLeave() {
		hoverX = null;
		draw();
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

<canvas
	bind:this={canvas}
	class="w-full flex-1 h-full"
	onpointermove={handlePointerMove}
	onpointerleave={handlePointerLeave}
></canvas>
