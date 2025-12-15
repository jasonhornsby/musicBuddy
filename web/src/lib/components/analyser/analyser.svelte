<script lang="ts">
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { Music, Info, HelpCircle, Eye } from 'lucide-svelte';
	import { VisualisationPanel } from './panels';

	type View = {
		id: string;
		name: string;
		description: string;
		icon: typeof Info;
	};

	const views: View[] = [
		{
			id: 'visualisation',
			name: 'Visualisation',
			description: 'Visual representation of audio data',
			icon: Eye
		},
	];

	let selectedView = $state<View>(views[0]);
</script>

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
		<header class="flex h-11 items-center border-b border-border/40 px-4">
			<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider">
				{selectedView.name}
			</span>
		</header>

		<!-- Content area -->
		<main class="flex-1 overflow-hidden bg-muted/10">
			{#if selectedView.id === 'visualisation'}
				<VisualisationPanel />
			{/if}
		</main>
	</div>
</div>
