<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Item from '$lib/components/ui/item';
	import { Button } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import { ChevronDown } from 'lucide-svelte';
	import { GeneralPanel, WaveformPanel } from './panels';

	type View = {
		id: string;
		name: string;
		description: string;
	};

	const views: View[] = [
		{ id: 'general', name: 'General Information', description: 'Overview of file metadata, duration, sample rate, and channels' },
		{ id: 'waveform', name: 'Waveform', description: 'Visual representation of amplitude over time' },
		{ id: 'spectral-flux', name: 'Spectral Flux', description: 'Measure of how quickly the spectrum changes, useful for onset detection' },
		{ id: 'spectrogram', name: 'Spectrogram', description: 'Frequency content visualization across time' }
	];

	let selectedView = $state<View>(views[0]);
</script>

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
		{/if}
	</main>
</div>
