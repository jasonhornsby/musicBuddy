<script lang="ts">
	import { EllipsisVertical, Workflow, Trash2 } from "lucide-svelte";
	import VisualisationCanvas from "../analyser/panels/visualisation.svelte";
	import PipelineGraph from "./pipeline-graph.svelte";
	import Button from "../ui/button/button.svelte";
	import ButtonGroup from "../ui/button-group/button-group.svelte";
	import {
		DropdownMenu,
		DropdownMenuTrigger,
		DropdownMenuContent,
		DropdownMenuItem,
	} from "../ui/dropdown-menu";
	import {
		Select,
		SelectTrigger,
		SelectContent,
		SelectItem,
	} from "../ui/select";
	import type { ChannelMode } from "$lib/types/nodeConfig";
	import { getAudioContext } from "$lib/context/audio.svelte";
	import { SpectrumVisualisation, WaveformVisualisation, Visualisation } from "$lib/utils/visualiser.svelte";


	const audioContext = getAudioContext();
	const workerManager = audioContext.getAudioWorkerManager();
	const waveformViz = new WaveformVisualisation(300);
	const spectrumViz = new SpectrumVisualisation(1000);

	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	const visualisations: Visualisation<any>[] = [waveformViz, spectrumViz];

	workerManager.createVisualizer(waveformViz);
	workerManager.createVisualizer(spectrumViz);

	const trackName = $state("Track 1");

	let channel = $state<ChannelMode>('mix');
	let selectedVizId = $state<string>(waveformViz.id);
	let showRenderGraph = $state(false);

	const selectedVisualisation = $derived(
		visualisations.find((v) => v.id === selectedVizId) ?? spectrumViz
	);

	$effect(() => {
		// waveformViz.updateConfig({ channel});
		spectrumViz.updateConfig({ channel, windowSize: 1024 });
	});
</script>

<div
	class="h-[280px] w-full rounded-lg border border-slate-200 bg-white flex flex-col md:flex-row overflow-hidden shadow-sm"
>
	<aside
		class="w-full md:w-48 md:max-w-[28%] border-r border-slate-200 bg-slate-50/80 text-[11px] flex flex-col"
	>
		<div class="flex items-center justify-between gap-1.5 px-2 py-2 border-b border-slate-100">
			<div class="flex flex-col gap-0.5 min-w-0">
				<h2 class="text-[12px] font-semibold text-slate-800 truncate">
					{trackName}
				</h2>
			</div>
			<DropdownMenu>
				<DropdownMenuTrigger>
					{#snippet child({ props })}
						<Button size="icon-sm" variant="ghost" class="h-7 w-7" {...props}>
							<EllipsisVertical class="w-3.5 h-3.5 text-slate-500" />
						</Button>
					{/snippet}
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuItem onclick={() => (showRenderGraph = !showRenderGraph)}>
						<Workflow class="w-4 h-4" />
						{showRenderGraph ? 'Hide render graph' : 'Show render graph'}
					</DropdownMenuItem>
					<DropdownMenuItem onclick={() => console.log('Remove track clicked')}>
						<Trash2 class="w-4 h-4" />
						Remove track
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>

		<div class="flex flex-col gap-2 px-2 py-2">
			<div class="flex flex-col gap-1">
				<span class="text-[10px] font-medium text-slate-500 uppercase tracking-[0.14em]">
					Channel
				</span>
				<ButtonGroup class="mt-0.5 w-full">
					<Button
						variant={channel === 'left' ? 'default' : 'outline'}
						size="sm"
						class="flex-1 text-[11px] px-1.5"
						onclick={() => (channel = 'left')}
					>
						Left
					</Button>
					<Button
						variant={channel === 'right' ? 'default' : 'outline'}
						size="sm"
						class="flex-1 text-[11px] px-1.5"
						onclick={() => (channel = 'right')}
					>
						Right
					</Button>
					<Button
						variant={channel === 'mix' ? 'default' : 'outline'}
						size="sm"
						class="flex-1 text-[11px] px-1.5"
						onclick={() => (channel = 'mix')}
					>
						Mix
					</Button>
				</ButtonGroup>
			</div>
			<div class="flex flex-col gap-1">
				<span class="text-[10px] font-medium text-slate-500 uppercase tracking-[0.14em]">
					Visualisation
				</span>
				<Select type="single" bind:value={selectedVizId}>
					<SelectTrigger class="w-full text-[11px] h-8">
						{selectedVisualisation.name}
					</SelectTrigger>
					<SelectContent align="start">
						{#each visualisations as viz}
							<SelectItem value={viz.id}>
								{#snippet children({ selected })}
									<div class="flex flex-col gap-0.5">
										<span class="font-medium">{viz.name}</span>
										<span class="text-[10px] text-muted-foreground">{viz.description}</span>
									</div>
								{/snippet}
							</SelectItem>
						{/each}
					</SelectContent>
				</Select>
			</div>
		</div>
	</aside>
	<div class="flex-[300px] h-[300px] md:h-auto md:flex-1">
		{#if showRenderGraph}
			<PipelineGraph data={workerManager.renderTree} />
		{:else if selectedVisualisation.ready}
			<VisualisationCanvas visualisation={selectedVisualisation} />
		{:else}
			<span>Loading...</span>
		{/if}
	</div>
</div>
