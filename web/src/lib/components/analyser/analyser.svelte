<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Item from '$lib/components/ui/item';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { Button } from '$lib/components/ui/button';
	import {
		ChevronDown,
		Music,
		Info,
		Activity,
		AudioWaveform,
		BarChart3,
		RotateCcw,
		HelpCircle
	} from 'lucide-svelte';
	import { GeneralPanel, SpectralFluxPanel, WaveformPanel } from './panels';
	import { getAudioContext } from '$lib/context/audio.svelte';

	const audioContext = getAudioContext();

	type View = {
		id: string;
		name: string;
		description: string;
		icon: typeof Info;
	};

	const views: View[] = [
		{
			id: 'waveform',
			name: 'Waveform',
			description: 'Visual representation of amplitude over time',
			icon: AudioWaveform
		},
		{
			id: 'general',
			name: 'General Information',
			description: 'Overview of file metadata, duration, sample rate, and channels',
			icon: Info
		},
		{
			id: 'spectral-flux',
			name: 'Spectral Flux',
			description: 'Measure of how quickly the spectrum changes, useful for onset detection',
			icon: Activity
		},
		{
			id: 'spectrogram',
			name: 'Spectrogram',
			description: 'Frequency content visualization across time',
			icon: BarChart3
		}
	];

	let selectedView = $state<View>(views[0]);
</script>

<Tooltip.Provider delayDuration={0}>
	<div class="flex h-dvh w-full bg-background">
		<!-- Figma-style icon sidebar -->
		<aside class="flex w-12 flex-col border-r border-border/40 bg-sidebar">
			<!-- Logo -->
			<div class="flex h-11 items-center justify-center border-b border-border/40">
				<Music class="h-4 w-4 text-foreground/70" strokeWidth={1.5} />
			</div>

			<!-- Tool icons -->
			<nav class="flex flex-1 flex-col items-center gap-1 py-2">
				{#each views as view (view.id)}
					<Tooltip.Root>
						<Tooltip.Trigger>
							<button
								class="flex h-8 w-8 items-center justify-center rounded-md transition-colors {selectedView.id ===
								view.id
									? 'bg-accent text-accent-foreground'
									: 'text-muted-foreground hover:bg-accent/50 hover:text-foreground'}"
								onclick={() => (selectedView = view)}
							>
								<view.icon class="h-4 w-4" strokeWidth={1.5} />
							</button>
						</Tooltip.Trigger>
						<Tooltip.Content side="right" sideOffset={8}>
							{view.name}
						</Tooltip.Content>
					</Tooltip.Root>
				{/each}
			</nav>

			<!-- Bottom actions -->
			<div class="flex flex-col items-center gap-1 border-t border-border/40 py-2">
				<Tooltip.Root>
					<Tooltip.Trigger>
						<button
							class="flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent/50 hover:text-foreground"
							onclick={() => audioContext.resetAudio()}
						>
							<RotateCcw class="h-4 w-4" strokeWidth={1.5} />
						</button>
					</Tooltip.Trigger>
					<Tooltip.Content side="right" sideOffset={8}>Reset</Tooltip.Content>
				</Tooltip.Root>

				<Tooltip.Root>
					<Tooltip.Trigger>
						<button
							class="flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent/50 hover:text-foreground"
						>
							<HelpCircle class="h-4 w-4" strokeWidth={1.5} />
						</button>
					</Tooltip.Trigger>
					<Tooltip.Content side="right" sideOffset={8}>Help</Tooltip.Content>
				</Tooltip.Root>
			</div>
		</aside>

		<!-- Main content area -->
		<div class="flex flex-1 flex-col min-w-0">
			<!-- Header bar -->
			<header class="flex h-11 items-center justify-between border-b border-border/40 px-4">
				<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider">
					{selectedView.name}
				</span>

				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Button
								variant="ghost"
								size="sm"
								class="h-7 px-2 text-xs font-normal text-muted-foreground hover:text-foreground"
								{...props}
							>
								View
								<ChevronDown class="ml-1 h-3 w-3" />
							</Button>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content class="w-72" align="end">
						{#each views as view (view.id)}
							<DropdownMenu.Item class="p-0" onclick={() => (selectedView = view)}>
								<Item.Root size="sm" class="w-full p-2">
									<div class="flex h-8 w-8 items-center justify-center rounded-md bg-muted/50">
										<view.icon class="h-4 w-4 text-muted-foreground" strokeWidth={1.5} />
									</div>
									<Item.Content class="gap-0.5">
										<Item.Title class="text-sm">{view.name}</Item.Title>
										<Item.Description class="text-xs">{view.description}</Item.Description>
									</Item.Content>
								</Item.Root>
							</DropdownMenu.Item>
						{/each}
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</header>

			<!-- Content area -->
			<main class="flex-1 overflow-hidden bg-muted/10">
				{#if selectedView.id === 'general'}
					<GeneralPanel />
				{:else if selectedView.id === 'waveform'}
					<WaveformPanel />
				{:else if selectedView.id === 'spectral-flux'}
					<SpectralFluxPanel />
				{:else}
					<div class="flex items-center justify-center h-full text-muted-foreground">
						<p class="text-sm">{selectedView.name} - Coming soon</p>
					</div>
				{/if}
			</main>
		</div>
	</div>
</Tooltip.Provider>

<!--

<div class="flex flex-col h-dvh">
	<header class="px-5 py-4 flex items-center justify-between">
		<h2 class="text-sm font-medium tracking-wide uppercase text-muted-foreground">
			{selectedView.name}
		</h2>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Button variant="outline" size="sm" {...props}>
						{selectedView.name}
						<ChevronDown class="ml-2 size-4" />
					</Button>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content class="w-80" align="end">
				{#each views as view (view.id)}
					<DropdownMenu.Item class="p-0" onclick={() => (selectedView = view)}>
						<Item.Root size="sm" class="w-full p-2">
							<Item.Content class="gap-0.5">
								<Item.Title>{view.name}</Item.Title>
								<Item.Description>{view.description}</Item.Description>
							</Item.Content>
						</Item.Root>
					</DropdownMenu.Item>
				{/each}
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</header>

	<Separator />

	<main class="flex-1 overflow-hidden">
		{#if selectedView.id === 'general'}
			<GeneralPanel />
		{:else if selectedView.id === 'waveform'}
			<WaveformPanel />
		{:else if selectedView.id === 'spectral-flux'}
			<SpectralFluxPanel />
		{/if}
	</main>
</div>


-->
