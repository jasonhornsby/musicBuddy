<script lang="ts">
	import { EllipsisVertical, Trash2 } from "lucide-svelte";
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
	import {
		ResizablePaneGroup,
		ResizablePane,
		ResizableHandle,
	} from "../ui/resizable";
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

	let selectedVizId = $state<string>(waveformViz.id);

	const selectedVisualisation = $derived(
		visualisations.find((v) => v.id === selectedVizId) ?? spectrumViz
	);

	const configSchema = $derived(workerManager.configSchemas.get(selectedVizId) ?? []);

	function getConfigValue(key: string): unknown {
		return selectedVisualisation.config[key as keyof typeof selectedVisualisation.config];
	}

	function updateConfig(key: string, value: unknown) {
		selectedVisualisation.updateConfig({ [key]: value });
	}
</script>

<div class="h-full w-full flex flex-row overflow-hidden">
	<!-- Sidebar -->
	<aside class="w-48 shrink-0 border-r border-slate-200 bg-slate-50/80 text-[11px] flex flex-col">
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
					<DropdownMenuItem onclick={() => console.log('Remove track clicked')}>
						<Trash2 class="w-4 h-4" />
						Remove track
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>

		<div class="flex flex-col gap-2 px-2 py-2 overflow-y-auto">
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

			{#each configSchema as param (param.key)}
				<div class="flex flex-col gap-1">
					<span class="text-[10px] font-medium text-slate-500 uppercase tracking-[0.14em]">
						{param.label}
					</span>
					{#if param.type === 'select' && param.options}
						<ButtonGroup class="mt-0.5 w-full">
							{#each param.options as option}
								<Button
									variant={getConfigValue(param.key) === option ? 'default' : 'outline'}
									size="sm"
									class="flex-1 text-[11px] px-1.5 capitalize"
									onclick={() => updateConfig(param.key, option)}
								>
									{option}
								</Button>
							{/each}
						</ButtonGroup>
					{:else if param.type === 'int' || param.type === 'float'}
						<div class="flex items-center gap-2">
							<input
								type="range"
								min={param.min ?? 0}
								max={param.max ?? 100}
								step={param.type === 'int' ? (param.step ?? 1) : (param.step ?? 0.1)}
								value={getConfigValue(param.key) as number}
								oninput={(e) => {
									const val =
										param.type === 'int'
											? parseInt(e.currentTarget.value)
											: parseFloat(e.currentTarget.value);
									updateConfig(param.key, val);
								}}
								class="flex-1 h-2 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-slate-700"
							/>
							<span class="text-[11px] text-slate-600 min-w-[3ch] text-right">
								{getConfigValue(param.key)}
							</span>
						</div>
					{/if}
				</div>
			{/each}
		</div>
	</aside>

	<!-- Main content area with resizable panels -->
	<ResizablePaneGroup direction="vertical" class="flex-1 min-w-0">
		<!-- Visualization area -->
		<ResizablePane defaultSize={70} minSize={30}>
			<div class="w-full h-full bg-white">
				{#if selectedVisualisation.ready}
					<VisualisationCanvas visualisation={selectedVisualisation} />
				{:else}
					<div class="w-full h-full flex items-center justify-center text-slate-400 text-sm">
						Loading...
					</div>
				{/if}
			</div>
		</ResizablePane>

		<ResizableHandle withHandle />

		<!-- Pipeline graph - always visible at the bottom -->
		<ResizablePane defaultSize={30} minSize={15}>
			<div class="w-full h-full bg-slate-900">
				<PipelineGraph data={workerManager.renderTree} />
			</div>
		</ResizablePane>
	</ResizablePaneGroup>
</div>
